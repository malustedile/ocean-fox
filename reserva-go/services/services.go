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
	
)
// --- Global Variables (Collections, Channels, etc.) ---
var (
	MongoClient          *mongo.Client
	DestinosCollection   *mongo.Collection
	ReservasCollection   *mongo.Collection
	InscricoesCollection   *mongo.Collection
	RabbitMQConnection   *amqp.Connection
	RabbitMQChannelGlobal *amqp.Channel // General purpose channel for publishing reserva-criada
)


// For RabbitMQ message content in pagamento-aprovado consumer
type PedidoReservaRMQ struct {
	ID string `json:"id"` 
	
}



// --- Helper Functions ---
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
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
	InscricoesCollection = db.Collection("inscricoes")

	RabbitMQConnection, err = amqp.Dial(rabbitMQURI)
	failOnError(err, "Failed to connect to RabbitMQ")
	log.Println("Connected to RabbitMQ!")

	// General channel for publishing to "reserva-criada-exc"
	RabbitMQChannelGlobal, err = RabbitMQConnection.Channel()
	failOnError(err, "Failed to open a global RabbitMQ channel")
}
