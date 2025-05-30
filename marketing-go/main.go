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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Inscricao struct {
    SessionId string    `bson:"sessionId" json:"sessionId"`
    Destino   string    `bson:"destino" json:"destino"`
    CriadoEm  time.Time `bson:"criadoEm" json:"criadoEm"`
}

type Promocao struct {
    SessionId string    `bson:"sessionId" json:"sessionId"`
    Mensagem  string    `bson:"mensagem" json:"mensagem"`
    CriadoEm  time.Time `bson:"criadoEm" json:"criadoEm"`
    Destino   string    `bson:"destino" json:"destino"`
}

func main() {
    ctx := context.Background()

    // MongoDB
    mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:exemplo123@localhost:27017"))
    if err != nil {
        log.Fatal(err)
    }
    db := mongoClient.Database("ocean-fox")
    inscricoes := db.Collection("inscricoes")
    promocoes := db.Collection("promocoes")

    // RabbitMQ
    rabbitConn, err := amqp.Dial("amqp://localhost")
    if err != nil {
        log.Fatal(err)
    }
    defer rabbitConn.Close()
    channel, err := rabbitConn.Channel()
    if err != nil {
        log.Fatal(err)
    }
    defer channel.Close()

    // HTTP Server
    app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowCredentials: true,
	}))
    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello Fiber")
    })

    app.Get("/minhas-inscricoes", func(c *fiber.Ctx) error {
        sessionId := c.Cookies("sessionId")
        var mySubscriptions []Inscricao
        var myPromotions []Promocao

        cursor, _ := inscricoes.Find(ctx, bson.M{"sessionId": sessionId})
    	cursor.All(ctx, &mySubscriptions)

        cursor2, _ := promocoes.Find(ctx, bson.M{"sessionId": sessionId})
        cursor2.All(ctx, &myPromotions)

        return c.JSON(fiber.Map{
            "subscriptions": mySubscriptions,
            "promotions":    myPromotions,
        })
    })

    app.Post("/inscrever", func(c *fiber.Ctx) error {
        sessionId := c.Cookies("sessionId")
        var body struct {
            Destino string `json:"destino"`
        }
        if err := c.BodyParser(&body); err != nil {
            return err
        }
        exchange := "promocoes-" + body.Destino
        channel.ExchangeDeclare(exchange, "fanout", false, false, false, false, nil)

        inscricoes.InsertOne(ctx, Inscricao{
            SessionId: sessionId,
            Destino:   body.Destino,
            CriadoEm:  time.Now(),
        })

        queue, _ := channel.QueueDeclare("", false, true, true, false, nil)
        channel.QueueBind(queue.Name, "", exchange, false, nil)

        msgs, _ := channel.Consume(queue.Name, "", true, true, false, false, nil)
        go func() {
            for d := range msgs {
                var promocao Promocao
                json.Unmarshal(d.Body, &promocao)
                promocao.SessionId = sessionId
                promocoes.InsertOne(ctx, promocao)
                fmt.Printf("Promoção recebida: %+v\n", promocao)
            }
        }()

        return c.JSON(fiber.Map{"success": true})
    })

    app.Post("/promocao", func(c *fiber.Ctx) error {
        var body struct {
            Destino  string `json:"destino"`
            Mensagem string `json:"mensagem"`
        }
        if err := c.BodyParser(&body); err != nil {
            return err
        }
        exchange := "promocoes-" + body.Destino
        channel.ExchangeDeclare(exchange, "fanout", false, false, false, false, nil)

        promocao := Promocao{
            Mensagem: body.Mensagem,
            Destino:  body.Destino,
            CriadoEm: time.Now(),
        }
        data, _ := json.Marshal(promocao)
        channel.Publish(exchange, "", false, false, amqp.Publishing{
            ContentType: "application/json",
            Body:        data,
        })
        fmt.Printf("[Publisher] Promoção enviada para exchange %s: %s\n", exchange, body.Mensagem)
        return nil
    })

    log.Fatal(app.Listen(":3004"))
}