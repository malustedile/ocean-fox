package services

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// --- Configuration Constants ---
const (
	mongoURI                = "mongodb://root:exemplo123@localhost:27017"
	rabbitMQURI             = "amqp://guest:guest@localhost:5672/"
	databaseName            = "ocean-fox"
	destinosCollectionName  = "destinos"
	reservasCollectionName  = "reservas"
	publicKeyPath           = "./pagamento.pub"
	httpServerAddress       = ":3000"
	reservaCriadaQueue      = "reserva-criada"
	ReservaCriadaExchange   = "reserva-criada-exc"    // Fanout
	reservaCanceladaQueue   = "reserva-cancelada"     // Direct
	ReservaCanceladaExchange = "reserva-cancelada-exc" // Direct
	pagamentoAprovadoExchange = "pagamento-aprovado-exc" // Direct
	pagamentoAprovadoRK     = "pagamento-aprovado"
	pagamentoRecusadoQueue  = "pagamento-recusado"
	pagamentoRecusadoExchange = "pagamento-recusado-exc" // Direct
	pagamentoRecusadoRK = "pagamento-recusado"
	bilheteGeradoQueue      = "bilhete-gerado"
)
// --- Global Variables (Collections, Channels, etc.) ---
var (
	MongoClient          *mongo.Client
	DestinosCollection   *mongo.Collection
	ReservasCollection   *mongo.Collection
	RabbitMQConnection   *amqp.Connection
	RabbitMQChannelGlobal *amqp.Channel // General purpose channel for publishing reserva-criada
)


// For RabbitMQ message content in pagamento-aprovado consumer
type PedidoReservaRMQ struct {
	ID string `json:"id"` 
	
}

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



func InitMongoDB() {
	var err error
	clientOptions := options.Client().ApplyURI(mongoURI)
	MongoClient, err = mongo.Connect(context.Background(), clientOptions)
	failOnError(err, "Failed to connect to MongoDB")

	err = MongoClient.Ping(context.Background(), nil)
	failOnError(err, "Failed to ping MongoDB")
	log.Println("Connected to MongoDB!")

	db := MongoClient.Database(databaseName)
	DestinosCollection = db.Collection(destinosCollectionName)
	ReservasCollection = db.Collection(reservasCollectionName)
}

func InitRabbitMQ() {
	var err error
	RabbitMQConnection, err = amqp.Dial(rabbitMQURI)
	failOnError(err, "Failed to connect to RabbitMQ")
	log.Println("Connected to RabbitMQ!")

	// General channel for publishing to "reserva-criada-exc"
	RabbitMQChannelGlobal, err = RabbitMQConnection.Channel()
	failOnError(err, "Failed to open a global RabbitMQ channel")
	
	// Declare queue "reserva-criada" (consumed by other services, durable)
	_, err = RabbitMQChannelGlobal.QueueDeclare(
		reservaCriadaQueue, // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	failOnError(err, "Failed to declare queue 'reserva-criada'")

	// Declare exchange "reserva-criada-exc" (fanout, not durable as per JS)
	err = RabbitMQChannelGlobal.ExchangeDeclare(
		ReservaCriadaExchange, // name
		"fanout",              // type
		false,                 // durable (JS has false)
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	failOnError(err, "Failed to declare exchange 'reserva-criada-exc'")

	// Channel for "reserva-cancelada"
	chReservaCancelada, err := RabbitMQConnection.Channel()
	failOnError(err, "Failed to open channel for reserva cancelada")
	defer chReservaCancelada.Close() // Close if only used for declaration here
	_, err = chReservaCancelada.QueueDeclare(
		reservaCanceladaQueue, // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	failOnError(err, "Failed to declare queue 'reserva-cancelada'")

	// Declare exchange "reserva-cancelada-exc" (fanout, not durable as per JS)
	err = RabbitMQChannelGlobal.ExchangeDeclare(
		ReservaCanceladaExchange, // name
		"fanout",              // type
		false,                 // durable (JS has false)
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	failOnError(err, "Failed to declare exchange 'reserva-cancelada-exc'")

	// Bind queue to exchange
	err = chReservaCancelada.QueueBind(
		reservaCanceladaQueue, // queue name
		"",                    // routing key (fanout ignores this)
		ReservaCanceladaExchange, // exchange
		false,                 // no-wait
		nil,                   // arguments
	)
	failOnError(err, "Failed to bind reserva cancelada queue to exchange")

	// Setup for Pagamento Aprovado consumer
	chPagamentoAprovado, err := RabbitMQConnection.Channel()
	failOnError(err, "Failed to open channel for pagamento aprovado")
	// Declare exchange "pagamento-aprovado-exc" (direct, durable)
	err = chPagamentoAprovado.ExchangeDeclare(
		pagamentoAprovadoExchange, // name
		"direct",                  // type
		true,                      // durable
		false,                     // auto-deleted
		false,                     // internal
		false,                     // no-wait
		nil,                       // arguments
	)
	failOnError(err, "Failed to declare exchange 'pagamento-aprovado-exc'")

	// Declare anonymous, durable queue for pagamento aprovado
	qPagamentoAprovado, err := chPagamentoAprovado.QueueDeclare(
		"",    // name (server-generated)
		true,  // durable
		false, // delete when unused
		true,  // exclusive (JS was true with empty name)
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare queue for pagamento aprovado")

	err = chPagamentoAprovado.QueueBind(
		qPagamentoAprovado.Name,   // queue name
		pagamentoAprovadoRK,       // routing key
		pagamentoAprovadoExchange, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind pagamento aprovado queue")
	
	// Start consumer for pagamento aprovado
	go consumePagamentoAprovado(chPagamentoAprovado, qPagamentoAprovado.Name)


	// Declare other queues (as per JS, these might be consumed by other services or this one later)
	// Channel for "pagamento-recusado"
	chPagamentoRecusado, err := RabbitMQConnection.Channel()
	failOnError(err, "Failed to open channel for pagamento recusado")
	// Declare exchange "pagamento-recusado-exc" (direct, durable)
	err = chPagamentoAprovado.ExchangeDeclare(
		pagamentoRecusadoExchange, // name
		"direct",                  // type
		true,                      // durable
		false,                     // auto-deleted
		false,                     // internal
		false,                     // no-wait
		nil,                       // arguments
	)
	failOnError(err, "Failed to declare exchange 'pagamento-recusado-exc'")

	// Declare anonymous, durable queue for pagamento aprovado
	qPagamentoRecusado, err := chPagamentoRecusado.QueueDeclare(
		"",    // name (server-generated)
		true,  // durable
		false, // delete when unused
		true,  // exclusive (JS was true with empty name)
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare queue for pagamento recusado")

	err = chPagamentoRecusado.QueueBind(
		qPagamentoRecusado.Name,   // queue name
		pagamentoRecusadoRK,       // routing key
		pagamentoRecusadoExchange, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind pagamento recusado queue")

	go consumePagamentoRecusado(chPagamentoRecusado, pagamentoRecusadoQueue)

	// Channel for "bilhete-gerado"
	chBilheteGerado, err := RabbitMQConnection.Channel()
	failOnError(err, "Failed to open channel for bilhete gerado")
	defer chBilheteGerado.Close() // Close if only used for declaration here
	_, err = chBilheteGerado.QueueDeclare(bilheteGeradoQueue, true, false, false, false, nil)
	failOnError(err, "Failed to declare queue 'bilhete-gerado'")
}
