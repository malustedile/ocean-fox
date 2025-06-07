package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"marketing/services"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)



func MySubscriptions(c *fiber.Ctx) error {
	ctx := context.Background()
	sessionId := c.Cookies("sessionId")
	var myPromotions []services.Promocao

	cursor2, _ := services.CollectionPromocoes.Find(ctx, bson.M{"sessionId": sessionId})
	cursor2.All(ctx, &myPromotions)

	return c.JSON(myPromotions)
}

func Subscribe(c *fiber.Ctx) error {
	ctx := context.Background()

	sessionId := c.Cookies("sessionId")
	exchange := "promocoes"
	services.RabbitMQChannel.ExchangeDeclare(exchange, "fanout", false, false, false, false, nil)

	services.CollectionInscricoes.InsertOne(ctx, services.Inscricao{
		SessionId: sessionId,
		CriadoEm:  time.Now(),
	})

	queue, _ := services.RabbitMQChannel.QueueDeclare("", false, true, true, false, nil)
	services.RabbitMQChannel.QueueBind(queue.Name, "", exchange, false, nil)

	msgs, _ := services.RabbitMQChannel.Consume(queue.Name, "", true, true, false, false, nil)

	stopChan := make(chan struct{})
	services.RabbitMQconsumers[sessionId] = stopChan

	go func() {
		for {
			select {
			case d := <-msgs:
				var promocao services.Promocao
				json.Unmarshal(d.Body, &promocao)
				promocao.SessionId = sessionId
				services.CollectionPromocoes.InsertOne(ctx, promocao)
				fmt.Printf("Promoção recebida: %+v\n", promocao)
			case <-stopChan:
				// Cancela o consumo e deleta a fila
					services.RabbitMQChannel.QueueDelete(queue.Name, false, false, false)
				return
			}
		}
	}()

	return c.JSON(fiber.Map{"success": true})
}

func CreatePromotion(c *fiber.Ctx) error {
	var body struct {
		Mensagem string `json:"mensagem"`
	}
	if err := c.BodyParser(&body); err != nil {
		return err
	}
	exchange := "promocoes"
	services.RabbitMQChannel.ExchangeDeclare(exchange, "fanout", false, false, false, false, nil)

	promocao := services.Promocao{
		Mensagem: body.Mensagem,
		CriadoEm: time.Now(),
	}
	data, _ := json.Marshal(promocao)
	services.RabbitMQChannel.Publish(exchange, "", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	})
	fmt.Printf("[Publisher] Promoção enviada para exchange %s: %s\n", exchange, body.Mensagem)
	return nil
}

func CancelSubscription(c *fiber.Ctx) error {
	ctx := context.Background()

	sessionId := c.Cookies("sessionId")
	if sessionId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "error": "sessionId não encontrado"})
	}

	res, err := services.CollectionInscricoes.DeleteOne(ctx, bson.M{"sessionId": sessionId})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "error": err.Error()})
	}
	if res.DeletedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "error": "Inscrição não encontrada"})
	}

	// Para o consumo da fila
	if stopChan, ok := services.RabbitMQconsumers[sessionId]; ok {
		close(stopChan)
		delete(services.RabbitMQconsumers, sessionId)
	}

	return c.JSON(fiber.Map{"success": true})
}