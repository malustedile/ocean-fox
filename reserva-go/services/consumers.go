package services

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --- Helper Functions ---
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func consumePagamentoAprovado(ch *amqp.Channel, queueName string) {
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack (false for manual ack)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer for pagamento aprovado")

	log.Printf(" [*] Waiting for 'pagamento-aprovado' messages. To exit press CTRL+C")
	for d := range msgs {
		log.Printf("Received a 'pagamento-aprovado' message: %s", d.Body)
		var payload PedidoPagamentoPayload
		if err := json.Unmarshal(d.Body, &payload); err != nil {
			log.Printf("Error unmarshalling 'pagamento-aprovado' message: %v", err)
			d.Nack(false, false) // Do not requeue
			continue
		}
		
		// Update MongoDB
		reservaID, err := primitive.ObjectIDFromHex(payload.Reserva.ID)
		if err != nil {
			log.Printf("Error converting reserva ID to ObjectID: %v", err)
			d.Nack(false, false) // Do not requeue
			continue
		}
		ReservasCollection.UpdateOne(
			context.Background(),
			bson.M{"_id": reservaID},
			bson.M{
				"$set": bson.M{
					"status":          "PAGAMENTO_APROVADO",
					"bilhete":         nil,
				},
			},
		)

		log.Printf("Reserva %s updated", reservaID.Hex())
		d.Ack(false) // Acknowledge the message
	}
}

func consumePagamentoRecusado(ch *amqp.Channel, queueName string) {
	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack (false for manual ack)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	failOnError(err, "Failed to register a consumer for pagamento recusado")

	log.Printf(" [*] Waiting for 'pagamento-recusado' messages. To exit press CTRL+C")
	for d := range msgs {
		log.Printf("Received a 'pagamento-recusado' message: %s", d.Body)
		var payload PedidoPagamentoPayload
		if err := json.Unmarshal(d.Body, &payload); err != nil {
			log.Printf("Error unmarshalling 'pagamento-recusado' message: %v", err)
			d.Nack(false, false) // Do not requeue
			continue
		}
		
		// Cancelando a reserva (publica como reserva cancelada)
		reservaID, err := primitive.ObjectIDFromHex(payload.Reserva.ID)
		if err != nil {
			log.Printf("Error converting reserva ID to ObjectID: %v", err)
			d.Nack(false, false) // Do not requeue
			continue
		}
		ReservasCollection.UpdateOne(
			context.Background(),
			bson.M{"_id": reservaID},
			bson.M{
				"$set": bson.M{
					"status":          "PAGAMENTO_RECUSADO",
					"bilhete":         nil,
				},
			},
		)
		// publicando na fila de reserva cancelada
		canceladaMsg := bson.M{
			"id": reservaID.Hex(),
			"destino":           payload.Reserva.Destino,
			"sessionId":         payload.Reserva.SessionID,
			"dataEmbarque":      payload.Reserva.DataEmbarque,
			"numeroPassageiros": payload.Reserva.NumeroPassageiros,
			"numeroCabines":     payload.Reserva.NumeroCabines,
			"valorTotal":        payload.Reserva.ValorTotal,
		}
		canceladaMsgBytes, err := json.Marshal(canceladaMsg)
		if err != nil {
			log.Printf("Error marshalling 'reserva cancelada' message: %v", err)
			d.Nack(false, false) // Do not requeue
			continue
		}
		err = RabbitMQChannelGlobal.PublishWithContext(
			context.Background(),
			ReservaCanceladaExchange, // exchange
			"",                    // routing key (fanout ignores this)
			false,                 // mandatory
			false,                 // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        canceladaMsgBytes,
			})
		if err != nil {
			log.Printf("Error publishing 'reserva cancelada' message: %v", err)
			d.Nack(false, false) // Do not requeue
			continue
		}

		log.Printf("Reserva %s cancelada devido a pagamento recusado", reservaID.Hex())
		d.Ack(false) // Acknowledge the message
	}
}
