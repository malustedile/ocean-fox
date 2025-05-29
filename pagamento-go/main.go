package main

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	amqp "github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReservaPayload struct {
    ID               string `json:"id"`
    Destino          string `json:"destino"`
    DataEmbarque     string `json:"dataEmbarque"`
    NumeroPassageiros int    `json:"numeroPassageiros"`
    NumeroCabines    int    `json:"numeroCabines"`
    LinkPagamento    string `json:"linkPagamento"`
    Status           string `json:"status"`
    CriadoEm         string `json:"criadoEm"`
}

type MensagemAssinada struct {
    Reserva    ReservaPayload `json:"reserva"`
    Assinatura string         `json:"assinatura"`
}

func main() {
    // MongoDB
    ctx := context.Background()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:exemplo123@localhost:27017"))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(ctx)
    db := client.Database("ocean-fox")
    reservas := db.Collection("reservas")

    // RabbitMQ
    rabbit, err := amqp.Dial("amqp://localhost")
    if err != nil {
        log.Fatal(err)
    }
    defer rabbit.Close()

    channelReserva, _ := rabbit.Channel()
    channelPagamentoAprovado, _ := rabbit.Channel()
    channelPagamentoRecusado, _ := rabbit.Channel()

    reservaExchange := "reserva-criada-exc"
    pagamentoAprovadoExchange := "pagamento-aprovado-exc"

    channelReserva.ExchangeDeclare(reservaExchange, "fanout", false, false, false, false, nil)
    channelReserva.QueueDeclare("reserva-criada", true, false, false, false, nil)
    channelPagamentoAprovado.QueueDeclare("pagamento-aprovado", true, false, false, false, nil)
    channelPagamentoAprovado.ExchangeDeclare(pagamentoAprovadoExchange, "direct", true, false, false, false, nil)
    channelPagamentoRecusado.QueueDeclare("pagamento-recusado", true, false, false, false, nil)
    channelReserva.QueueBind("reserva-criada", "", reservaExchange, false, nil)

    msgs, err := channelReserva.Consume("reserva-criada", "", false, false, false, false, nil)
    if err != nil {
        log.Fatal(err)
    }

    // Carrega chave privada
    privKeyData, err := ioutil.ReadFile("./private.key")
    if err != nil {
        log.Fatal("Erro ao ler chave privada:", err)
    }
    block, _ := pem.Decode(privKeyData)
    if block == nil {
        log.Fatal("Falha ao decodificar PEM da chave privada")
    }
    privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
    if err != nil {
        log.Fatal("Erro ao parsear chave privada:", err)
    }

    go func() {
        for d := range msgs {
            var reserva ReservaPayload
            if err := json.Unmarshal(d.Body, &reserva); err != nil {
                log.Println("Erro ao decodificar reserva:", err)
                d.Ack(false)
                continue
            }
            fmt.Println("Reserva recebida:", reserva)

            pagamentoAprovado := randBool()
            if pagamentoAprovado {
                reserva.Status = "PAGAMENTO_APROVADO"
            } else {
                reserva.Status = "PAGAMENTO_REPROVADO"
            }

            // Assinatura digital
            reservaBytes, _ := json.Marshal(reserva)
            hash := sha256.Sum256(reservaBytes)
            signature, err := rsa.SignPKCS1v15(rand.Reader, privKey.(*rsa.PrivateKey), crypto.SHA256, hash[:])
            if err != nil {
                log.Println("Erro ao assinar:", err)
                d.Ack(false)
                continue
            }
            assinatura := base64.StdEncoding.EncodeToString(signature)

            payload, _ := json.Marshal(MensagemAssinada{
                Reserva:    reserva,
                Assinatura: assinatura,
            })

            if pagamentoAprovado {
                channelPagamentoAprovado.Publish(
                    pagamentoAprovadoExchange,
                    "pagamento-aprovado",
                    false, false,
                    amqp.Publishing{
                        ContentType: "application/json",
                        Body:        payload,
                    },
                )
                fmt.Println("Pagamento aprovado:", reserva)
            } else {
                channelPagamentoRecusado.Publish(
                    "",
                    "pagamento-recusado",
                    false, false,
                    amqp.Publishing{
                        ContentType: "application/json",
                        Body:        payload,
                    },
                )
                fmt.Println("Pagamento recusado:", reserva)
            }

            // Atualiza no MongoDB
            objID, err := primitive.ObjectIDFromHex(reserva.ID)
            if err == nil {
                update := bson.M{
                    "$set": bson.M{
                        "status":    reserva.Status,
                        "assinatura": assinatura,
                    },
                }
                result, err := reservas.UpdateOne(ctx, bson.M{"_id": objID}, update)
                if err != nil {
                    log.Println("Erro ao atualizar reserva:", err)
                } else {
                    fmt.Println("MongoDB update:", result)
                }
            }

            d.Ack(false)
        }
    }()

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    })
    port := "3001"
    fmt.Printf("App is running at 0.0.0.0:%s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

func randBool() bool {
    return time.Now().UnixNano()%2 == 0
}