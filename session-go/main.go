package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Session struct {
    ID        primitive.ObjectID     `bson:"_id,omitempty" json:"_id,omitempty"`
    CreatedAt time.Time              `bson:"createdAt" json:"createdAt"`
    Data      map[string]interface{} `bson:"data" json:"data"`
}

func main() {
    // MongoDB
    mongoURI := "mongodb://root:exemplo123@localhost:27017"
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
    if err != nil {
        log.Fatal(err)
    }
    db := client.Database("ocean-fox")
    sessions := db.Collection("sessions")
    // reservas := db.Collection("reservas") // Não usado diretamente aqui

    // RabbitMQ
    rabbitConn, err := amqp.Dial("amqp://localhost")
    if err != nil {
        log.Fatal(err)
    }
    defer rabbitConn.Close()

    ch, err := rabbitConn.Channel()
    if err != nil {
        log.Fatal(err)
    }
    defer ch.Close()

    reservaExchange := "reserva-criada-exc"
    queue, err := ch.QueueDeclare(
        "reserva-criada-session",
        true,  // durable
        false, // delete when unused
        false, // exclusive
        false, // no-wait
        nil,   // arguments
    )
    if err != nil {
        log.Fatal(err)
    }

    err = ch.ExchangeDeclare(
        reservaExchange,
        "fanout",
        false, // durable
        false, // auto-deleted
        false, // internal
        false, // no-wait
        nil,   // arguments
    )
    if err != nil {
        log.Fatal(err)
    }

    err = ch.QueueBind(
        queue.Name,
        "",
        reservaExchange,
        false,
        nil,
    )
    if err != nil {
        log.Fatal(err)
    }

    msgs, err := ch.Consume(
        queue.Name,
        "",
        false, // auto-ack
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        log.Fatal(err)
    }

    go func() {
        for d := range msgs {
            var reserva map[string]interface{}
            if err := json.Unmarshal(d.Body, &reserva); err == nil {
                fmt.Printf("Reserva recebida: %+v\n", reserva)
            }
            d.Ack(false)
        }
    }()

    // Fiber app
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000", // coloque aqui o domínio do seu frontend
		AllowCredentials: true,
	}))

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("OK")
    })

    app.Get("/session", func(c *fiber.Ctx) error {
        cookie := c.Cookies("sessionId")
        var sessionData Session
        var sessionId string

        if cookie != "" {
            objID, err := primitive.ObjectIDFromHex(cookie)
            if err == nil {
                err = sessions.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&sessionData)
                if err == nil {
                    sessionId = cookie
                }
            }
        }

        if sessionId == "" {
            // Cria nova sessão
            sessionData = Session{
                CreatedAt: time.Now(),
                Data:      map[string]interface{}{},
            }
            res, err := sessions.InsertOne(context.TODO(), sessionData)
            if err != nil {
                return c.Status(500).JSON(fiber.Map{"erro": "Erro ao criar sessão"})
            }
            sessionId = res.InsertedID.(primitive.ObjectID).Hex()
            c.Cookie(&fiber.Cookie{
                Name:     "sessionId",
                Value:    sessionId,
                HTTPOnly: true,
                Path:     "/",
                Domain:   "localhost",
                SameSite: "None",
                Secure:   true,
                MaxAge:   60 * 60 * 24 * 365,
            })
        }

        return c.JSON(fiber.Map{
            "mensagem":   "Sessão ativa",
            "sessionId":  sessionId,
            "sessionData": sessionData,
        })
    })

    log.Fatal(app.Listen(":3005"))
}