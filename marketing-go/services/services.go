package services

// MongoDB
import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient *mongo.Client
	Db          *mongo.Database
	CollectionInscricoes  *mongo.Collection
	CollectionPromocoes   *mongo.Collection
	RabbitMQChannel *amqp.Channel
	RabbitMQConn *amqp.Connection
	ctx         = context.Background()
)

type Inscricao struct {
    SessionId string    `bson:"sessionId" json:"sessionId"`
    CriadoEm  time.Time `bson:"criadoEm" json:"criadoEm"`
}

type Promocao struct {
    SessionId string    `bson:"sessionId" json:"sessionId"`
    Mensagem  string    `bson:"mensagem" json:"mensagem"`
    CriadoEm  time.Time `bson:"criadoEm" json:"criadoEm"`
}
var RabbitMQconsumers = make(map[string]chan struct{})


func Init() {
	var err error
	MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:exemplo123@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	Db = MongoClient.Database("ocean-fox")
	CollectionInscricoes = Db.Collection("inscricoes")
	CollectionPromocoes = Db.Collection("promocoes")

	// RabbitMQ
    RabbitMQConn, err := amqp.Dial("amqp://localhost")
    if err != nil {
        log.Fatal(err)
    }
    RabbitMQChannel, err := RabbitMQConn.Channel()
    if err != nil {
        log.Fatal(err)
    }
    defer RabbitMQChannel.Close()
}

