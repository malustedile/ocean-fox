package consumers

import (
	"context"
	"encoding/json"
	"log"
	"reserva-go/handlers" // Make sure this path matches your go.mod module name and folder structure
	"reserva-go/services"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
)

type Promocao struct {
    SessionId string    `bson:"sessionId" json:"sessionId"`
    Mensagem  string    `bson:"mensagem" json:"mensagem"`
    CriadoEm  time.Time `bson:"criadoEm" json:"criadoEm"`
}

// Removed duplicate SSEMessage type; use handlers.SSEMessage instead.

type PedidoPagamentoPayload struct {
	Reserva struct {
		ID                string  `json:"id"` // This is the _id of the reserva
		Destino           string  `json:"destino"`
		SessionID         string  `json:"sessionId"`
		DataEmbarque      string  `json:"dataEmbarque"`
		NumeroPassageiros int     `json:"numeroPassageiros"`
		NumeroCabines     int     `json:"numeroCabines"`
		ValorTotal        float64 `json:"valorTotal"`
		LinkPagamento     string  `json:"linkPagamento"`
		Status            string  `json:"status"`
		// Bilhete might not be part of the signed data
		CriadoEm string `json:"criadoEm"` // ISO String from JS
	} `json:"reserva"`
	Assinatura string `json:"assinatura"`
}

// --- Helper Functions ---
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func consumePromocoes(ch *amqp.Channel, queueName string) {
    msgs, err := ch.Consume(
        queueName, // queue
        "",        // consumer
        false,     // auto-ack (false para ack manual)
        false,     // exclusive
        false,     // no-local
        false,     // no-wait
        nil,       // args
    )
    failOnError(err, "Failed to register a consumer for promocoes")

    log.Printf(" [*] Waiting for 'promocoes' messages. To exit press CTRL+C")
    for d := range msgs {
        log.Printf("Received a 'promocao' message: %s", d.Body)
        var promocao Promocao
        if err := json.Unmarshal(d.Body, &promocao); err != nil {
            log.Printf("Error unmarshalling 'promocao' message: %v", err)
            d.Nack(false, false) // Não reencaminha
            continue
        }

		cursor, err := services.InscricoesCollection.Find(context.Background(), bson.M{})
		if err != nil {
			log.Printf("Erro ao buscar inscritos: %v", err)
			d.Nack(false, false)
			continue
		}

		for cursor.Next(context.Background()) {
			var inscricao struct {
				SessionID string `bson:"sessionId"`
			}
			if err := cursor.Decode(&inscricao); err != nil {
				log.Printf("Erro ao decodificar inscricao: %v", err)
				continue
			}
			log.Printf("Sending to sessionId: %s", inscricao.SessionID)

			sseMsg := handlers.SSEMessage{
				SessionID: inscricao.SessionID,
				Msg:       promocao.Mensagem,
				EventType: "promocao",
			}
			handlers.SendMessageToClient(sseMsg)
		}
		if err := cursor.Err(); err != nil {
			log.Printf("Erro no cursor de inscritos: %v", err)
		}
		cursor.Close(context.Background())

	}
}
