package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"log"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Structs para os payloads (equivalente às interfaces TypeScript)
type ReservaPayload struct {
	Destino           string `json:"destino"`
	DataEmbarque      string `json:"dataEmbarque"`
	NumeroPassageiros int    `json:"numeroPassageiros"`
	NumeroCabines     int    `json:"numeroCabines"`
	LinkPagamento     string `json:"linkPagamento"`
	Status            string `json:"status"`
	CriadoEm          string `json:"criadoEm"`
}

type PedidoPayload struct {
	Reserva    ReservaPayload `json:"reserva"`
	Assinatura string         `json:"assinatura"` // Assinatura em Base64
}

type Bilhete struct {
	IDReserva         string `json:"idReserva"`
	Destino           string `json:"destino"`
	DataEmbarque      string `json:"dataEmbarque"`
	NumeroPassageiros int    `json:"numeroPassageiros"`
	NumeroCabines     int    `json:"numeroCabines"`
	CriadoEm          string `json:"criadoEm"`
}

const (
	rabbitMQURL               = "amqp://guest:guest@localhost:5672/" // Ajuste conforme necessário
	pagamentoAprovadoExchange = "pagamento-aprovado-exc"
	pagamentoAprovadoRoutingKey = "pagamento-aprovado"
	bilheteGeradoQueue        = "bilhete-gerado"
	publicKeyPath             = "./pagamento.pub" // Caminho para sua chave pública
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// --- Conexão com RabbitMQ ---
	conn, err := amqp.Dial(rabbitMQURL)
	failOnError(err, "Falha ao conectar ao RabbitMQ")
	defer conn.Close()

	chPagamentoAprovado, err := conn.Channel()
	failOnError(err, "Falha ao abrir canal para pagamento aprovado")
	defer chPagamentoAprovado.Close()

	chBilheteGerado, err := conn.Channel()
	failOnError(err, "Falha ao abrir canal para bilhete gerado")
	defer chBilheteGerado.Close()

	// --- Declaração de Exchange e Queues ---
	err = chPagamentoAprovado.ExchangeDeclare(
		pagamentoAprovadoExchange, // name
		"direct",                  // type
		true,                      // durable
		false,                     // auto-deleted
		false,                     // internal
		false,                     // no-wait
		nil,                       // arguments
	)
	failOnError(err, "Falha ao declarar exchange 'pagamento-aprovado-exc'")

	// Declara uma queue anônima e durável para o consumidor
	qPagamento, err := chPagamentoAprovado.QueueDeclare(
		"",    // name (vazio para nome gerado pelo servidor)
		true,  // durable
		false, // delete when unused
		true,  // exclusive (se true, só esta conexão pode usar. Se false, pode ser compartilhada)
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Falha ao declarar queue de pagamento aprovado")

	err = chPagamentoAprovado.QueueBind(
		qPagamento.Name,             // queue name
		pagamentoAprovadoRoutingKey, // routing key
		pagamentoAprovadoExchange,   // exchange
		false,
		nil,
	)
	failOnError(err, "Falha ao fazer bind da queue de pagamento aprovado")

	_, err = chBilheteGerado.QueueDeclare(
		bilheteGeradoQueue, // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	failOnError(err, "Falha ao declarar queue 'bilhete-gerado'")

	// --- Consumidor de Mensagens ---
	msgs, err := chPagamentoAprovado.Consume(
		qPagamento.Name, // queue
		"",              // consumer
		false,           // auto-ack (false pois vamos fazer manualmente)
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	failOnError(err, "Falha ao registrar consumidor")

	// Carrega a chave pública uma vez
	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	failOnError(err, "Falha ao ler arquivo da chave pública")

	pemBlock, _ := pem.Decode(publicKeyBytes)
	if pemBlock == nil {
		log.Fatalf("Falha ao decodificar bloco PEM da chave pública")
	}
	pub, err := x509.ParsePKIXPublicKey(pemBlock.Bytes)
	failOnError(err, "Falha ao parsear chave pública PKIX")

	rsaPubKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		log.Fatalf("Chave pública não é do tipo RSA")
	}

	var forever chan struct{} // Canal para manter o main rodando

	go func() {
		for msg := range msgs {
			log.Printf("Recebida mensagem: %s", msg.Body)
			var pedido PedidoPayload
			err := json.Unmarshal(msg.Body, &pedido)
			if err != nil {
				log.Printf("Erro ao fazer unmarshal do JSON do pedido: %s", err)
				msg.Nack(false, false) // Descarta a mensagem se não puder ser parseada
				continue
			}

			// Verifica a assinatura
			// 1. Serializa a reserva para JSON (da mesma forma que foi assinada)
			reservaJSON, err := json.Marshal(pedido.Reserva)
			if err != nil {
				log.Printf("Erro ao fazer marshal da reserva para verificação: %s", err)
				msg.Nack(false, false)
				continue
			}

			// 2. Decodifica a assinatura de Base64
			assinaturaBytes, err := base64.StdEncoding.DecodeString(pedido.Assinatura)
			if err != nil {
				log.Printf("Erro ao decodificar assinatura Base64: %s", err)
				msg.Nack(false, false)
				continue
			}

			// 3. Calcula o hash SHA256 da reserva serializada
			hasher := sha256.New()
			hasher.Write(reservaJSON)
			hashed := hasher.Sum(nil)

			// 4. Verifica a assinatura
			err = rsa.VerifyPKCS1v15(rsaPubKey, crypto.SHA256, hashed, assinaturaBytes)
			if err == nil {
				log.Println("Assinatura válida!")

				bilhete := Bilhete{
					IDReserva:         pedido.Reserva.LinkPagamento, // Usando LinkPagamento como ID, ajuste se necessário
					Destino:           pedido.Reserva.Destino,
					DataEmbarque:      pedido.Reserva.DataEmbarque,
					NumeroPassageiros: pedido.Reserva.NumeroPassageiros,
					NumeroCabines:     pedido.Reserva.NumeroCabines,
					CriadoEm:          time.Now().UTC().Format(time.RFC3339Nano),
				}

				bilheteJSON, err := json.Marshal(bilhete)
				if err != nil {
					log.Printf("Erro ao fazer marshal do bilhete JSON: %s", err)
					msg.Nack(false, true) // Nack e requeue, pode ser um erro transitório
					continue
				}

				err = chBilheteGerado.Publish(
					"",                 // exchange (default)
					bilheteGeradoQueue, // routing key (nome da queue quando exchange é default)
					false,              // mandatory
					false,              // immediate
					amqp.Publishing{
						ContentType:  "application/json",
						Body:         bilheteJSON,
						DeliveryMode: amqp.Persistent, // Torna a mensagem persistente
					})

				if err != nil {
					log.Printf("Erro ao publicar mensagem de bilhete gerado: %s", err)
					msg.Nack(false, true) // Nack e requeue
				} else {
					log.Printf("Bilhete gerado e enviado: %+v", bilhete)
					msg.Ack(false) // Confirma o processamento da mensagem original
				}

			} else {
				log.Printf("Assinatura inválida: %s", err)
				msg.Nack(false, false) // Descarta a mensagem
			}
		}
	}()

	// --- Servidor HTTP Simples (Equivalente ao Elysia) ---
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Olá do servidor Go!"))
	})

	serverAddr := "0.0.0.0:3002"
	log.Printf("🚀 Servidor Go rodando em %s", serverAddr)
	go func() {
		if err := http.ListenAndServe(serverAddr, nil); err != nil {
			log.Fatalf("Falha ao iniciar servidor HTTP: %s", err)
		}
	}()


	log.Printf(" [*] Aguardando mensagens. Para sair pressione CTRL+C")
	<-forever // Bloqueia a execução para o consumidor continuar rodando
}