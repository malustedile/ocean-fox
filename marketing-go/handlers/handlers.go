package handlers

import (
	"encoding/json"
	"fmt"
	"marketing/services"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
)

// Função auxiliar para obter sessionId e validar
func getSessionId(c *fiber.Ctx) (string, error) {
    sessionId := c.Cookies("sessionId")
    if sessionId == "" {
        return "", fmt.Errorf("sessionId não encontrado")
    }
    return sessionId, nil
}
func MySubscriptions(c *fiber.Ctx) error {
    ctx := c.Context()
    sessionId, err := getSessionId(c)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"hasSubscription": false, "promotions": []services.Promocao{}})
    }
    var myPromotions []services.Promocao

    var inscricao services.Inscricao
    err = services.CollectionInscricoes.FindOne(ctx, bson.M{"sessionId": sessionId}).Decode(&inscricao)
    if err != nil {
        return c.JSON(fiber.Map{
            "hasSubscription": false,
            "promotions":      []services.Promocao{},
        })
    }

    cursor, err := services.CollectionPromocoes.Find(ctx, bson.M{
        "criadoEm": bson.M{"$gt": inscricao.CriadoEm},
    })
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"hasSubscription": true, "promotions": []services.Promocao{}})
    }
    defer cursor.Close(ctx)
    if err := cursor.All(ctx, &myPromotions); err != nil {
        return c.Status(500).JSON(fiber.Map{"hasSubscription": true, "promotions": []services.Promocao{}})
    }

    // Ordena as promoções por CriadoEm decrescente
    sort.Slice(myPromotions, func(i, j int) bool {
        return myPromotions[i].CriadoEm.After(myPromotions[j].CriadoEm)
    })

    return c.JSON(fiber.Map{
        "hasSubscription": true,
        "promotions":      myPromotions,
    })
}

func Subscribe(c *fiber.Ctx) error {
    ctx := c.Context()
    sessionId, err := getSessionId(c)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
    }

    count, err := services.CollectionInscricoes.CountDocuments(ctx, bson.M{"sessionId": sessionId})
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"success": false, "error": "Erro ao verificar inscrição"})
    }
    if count > 0 {
        return c.JSON(fiber.Map{"success": false, "message": "Já inscrito"})
    }

    _, err = services.CollectionInscricoes.InsertOne(ctx, services.Inscricao{
        SessionId: sessionId,
        CriadoEm:  time.Now(),
    })
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"success": false, "error": "Erro ao inscrever"})
    }

    return c.JSON(fiber.Map{"success": true})
}

func CreatePromotion(c *fiber.Ctx) error {
    ctx := c.Context()

    var body struct {
        Mensagem string `json:"mensagem"`
    }
    if err := c.BodyParser(&body); err != nil || body.Mensagem == "" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Mensagem inválida"})
    }

    if services.RabbitMQChannel == nil {
        return c.Status(500).JSON(fiber.Map{"error": "RabbitMQ channel not initialized"})
    }

    promocao := services.Promocao{
        Mensagem: body.Mensagem,
        CriadoEm: time.Now(),
    }

    _, err := services.CollectionPromocoes.InsertOne(ctx, promocao)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Erro ao salvar promoção"})
    }

    data, _ := json.Marshal(promocao)

    err = services.RabbitMQChannel.Publish(
        "",          // exchange vazio para fila simples
        "promocoes", // routing key = nome da fila
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        data,
        },
    )
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    fmt.Printf("[Publisher] Promoção enviada para fila promocoes: %s\n", body.Mensagem)
    return c.JSON(fiber.Map{"success": true})
}

func CancelSubscription(c *fiber.Ctx) error {
    ctx := c.Context()
    sessionId, err := getSessionId(c)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": err.Error()})
    }

    res, err := services.CollectionInscricoes.DeleteMany(ctx, bson.M{"sessionId": sessionId})
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
    }
    if res.DeletedCount == 0 {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "Inscrição não encontrada"})
    }

    if stopChan, ok := services.RabbitMQconsumers[sessionId]; ok {
        close(stopChan)
        delete(services.RabbitMQconsumers, sessionId)
    }

    return c.JSON(fiber.Map{"success": true})
}