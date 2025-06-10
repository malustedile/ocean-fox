package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB
var db *mongo.Database

// Estruturas
type Destino struct {
	ID        string      `bson:"_id,omitempty" json:"id,omitempty"`
	Nome      string      `bson:"nome" json:"nome"`
	Categoria string      `bson:"categoria" json:"categoria"`
	Cabines   int         `bson:"cabines" json:"cabines"`
	Descricao DestinoInfo `bson:"descricao" json:"descricao"`
}

type DestinoInfo struct {
	DatasDisponiveis  []string `bson:"datasDisponiveis" json:"datasDisponiveis"`
	Navio             string   `bson:"navio" json:"navio"`
	Embarque          string   `bson:"embarque" json:"embarque"`
	Desembarque       string   `bson:"desembarque" json:"desembarque"`
	LugaresVisitados  []string `bson:"lugaresVisitados" json:"lugaresVisitados"`
	Noites            int      `bson:"noites" json:"noites"`
	ValorPorPessoa    int      `bson:"valorPorPessoa" json:"valorPorPessoa"`
}

// Simulação da estrutura das mensagens de reserva
type MensagemReserva struct {
	Destino   		string `json:"destino"`
	NumeroCabines   int    `json:"numeroCabines"`
}

func main() {

	ctx := context.TODO()

	// Conexão MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:exemplo123@localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	db = client.Database("ocean-fox")



	// Conexão RabbitMQ (compartilhada apenas a conexão, não o canal!)
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Erro ao conectar RabbitMQ:", err)
	}

	// Goroutine para consumir reserva criada
	go func() {
		ch, err := conn.Channel()
		if err != nil {
			log.Fatal("Erro ao abrir canal para reserva-criada:", err)
		}
		defer ch.Close()

		consumirExchange(ch, "reserva-criada-exc", false)
	}()

	// Goroutine para consumir reserva cancelada
	go func() {
		ch, err := conn.Channel()
		if err != nil {
			log.Fatal("Erro ao abrir canal para reserva-cancelada:", err)
		}
		defer ch.Close()

		consumirExchange(ch, "reserva-cancelada-exc", true)
	}()
		// Router mux
	r := mux.NewRouter()

	// Rotas com método explícito
	r.HandleFunc("/destinos", listarDestinosHandler).Methods("GET")
	r.HandleFunc("/destinos", criarDestinoHandler).Methods("POST")

	log.Println("Servidor iniciado em :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Handler GET /destinos
func listarDestinosHandler(w http.ResponseWriter, r *http.Request) {
	cursor, err := db.Collection("destinos").Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Erro ao buscar destinos", http.StatusInternalServerError)
		return
	}
	var destinos []Destino
	if err := cursor.All(context.TODO(), &destinos); err != nil {
		http.Error(w, "Erro ao processar destinos", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(destinos)
}

// Handler POST /destinos
func criarDestinoHandler(w http.ResponseWriter, r *http.Request) {
	var destino Destino
	if err := json.NewDecoder(r.Body).Decode(&destino); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	_, err := db.Collection("destinos").InsertOne(context.TODO(), destino)
	if err != nil {
		http.Error(w, "Erro ao inserir destino", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"mensagem": "Destino inserido com sucesso"})
}

func consumirExchange(ch *amqp.Channel, exchange string, cancelar bool) {
	// Declara a exchange (caso ainda não tenha sido declarada)
	err := ch.ExchangeDeclare(
		exchange,
		"fanout", // ou "direct"/"topic", dependendo do seu caso
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Erro ao declarar exchange %s: %v", exchange, err)
	}

	// Cria uma queue exclusiva (pode ser nomeada se quiser manter)
	q, err := ch.QueueDeclare(
		"",    // nome vazio cria uma queue exclusiva temporária
		false, // durável
		true,  // auto delete
		true,  // exclusiva
		false, // noWait
		nil,
	)
	if err != nil {
		log.Fatalf("Erro ao declarar fila para exchange %s: %v", exchange, err)
	}

	// Faz bind da fila à exchange
	err = ch.QueueBind(
		q.Name,
		"",        // routing key ("" para fanout, ou específico para direct/topic)
		exchange,  // nome da exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Erro ao fazer bind da fila %s à exchange %s: %v", q.Name, exchange, err)
	}

	// Começa a consumir da fila ligada à exchange
	msgs, err := ch.Consume(
		q.Name, "", true, false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("Erro ao consumir da fila ligada à exchange %s: %v", exchange, err)
	}

	for msg := range msgs {
		log.Printf("Mensagem recebida da exchange %s: %s", exchange, msg.Body)
		var res MensagemReserva
		if err := json.Unmarshal(msg.Body, &res); err != nil {
			log.Println("Erro ao parsear mensagem:", err)
			continue
		}
		qtd := res.NumeroCabines
		if cancelar {
			qtd = -qtd
		}
		filter := bson.M{"nome": res.Destino}
		update := bson.M{"$inc": bson.M{"cabines": -qtd}}
		_, err := db.Collection("destinos").UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Printf("Erro ao atualizar destino %s: %v", res.Destino, err)
		} else {
			log.Printf("Destino %s atualizado (exchange %s): %d cabines", res.Destino, exchange, -qtd)
		}
	}
}