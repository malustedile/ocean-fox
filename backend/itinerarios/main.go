package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// Categorias (equivalent to enum)
type Categoria string

const (
	CategoriaBrasil         Categoria = "Brasil"
	CategoriaAmericaDoSul   Categoria = "América do Sul"
	CategoriaCaribe         Categoria = "Caribe"
	CategoriaAmericaDoNorte Categoria = "América do Norte"
	CategoriaAfrica         Categoria = "África"
	CategoriaOrienteMedio   Categoria = "Oriente Médio"
	CategoriaAsia           Categoria = "Ásia"
	CategoriaMediterraneo   Categoria = "Mediterrâneo"
	CategoriaEscandinavia   Categoria = "Escandinávia"
	CategoriaOceania        Categoria = "Oceania"
)

var AllCategorias = []Categoria{
	CategoriaBrasil, CategoriaAmericaDoSul, CategoriaCaribe, CategoriaAmericaDoNorte,
	CategoriaAfrica, CategoriaOrienteMedio, CategoriaAsia, CategoriaMediterraneo,
	CategoriaEscandinavia, CategoriaOceania,
}

type FiltrosDTO struct {
	Destino   *string    `json:"destino"`
	Mes       *string    `json:"mes"`
	Embarque  *string    `json:"embarque"`
	Categoria *string `json:"categoria"`
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
	r.HandleFunc("/destinos-por-categoria", destinosPorCategoriaHandler).Methods("GET")
	r.HandleFunc("/destinos/buscar", BuscarDestinosHandler).Methods("POST")
	r.HandleFunc("/destinos", criarDestinoHandler).Methods("POST")

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173"}), // Or specify your frontend domain(s)
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}),
		handlers.AllowCredentials(), // Important if you use cookies/auth headers from frontend
	)

	log.Println("Servidor iniciado em :8080")
	if err := http.ListenAndServe(":8080", corsHandler(r)); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}


func BuscarDestinosHandler(w http.ResponseWriter, r *http.Request) {
	var dto FiltrosDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido.")
		return
	}

	filter := bson.M{}
	if dto.Destino != nil && *dto.Destino != "" {
		filter["descricao.lugaresVisitados"] = bson.M{
			"$elemMatch": bson.M{"$regex": primitive.Regex{Pattern: *dto.Destino, Options: "i"}},
		}
	}
	if dto.Embarque != nil && *dto.Embarque != "" {
		filter["descricao.embarque"] = bson.M{"$regex": primitive.Regex{Pattern: *dto.Embarque, Options: "i"}}
	}
	if dto.Mes != nil && *dto.Mes != "" {
		mesNum, err := strconv.Atoi(*dto.Mes)
		if err == nil && mesNum >= 1 && mesNum <= 12 {
			// Regex for YYYY-MM-DD or similar, matching the month part.
			// Example: "-06-" for June. This regex is specific to date formats "YYYY-MM-DD".
			monthPattern := fmt.Sprintf("-%02d-", mesNum)
			filter["descricao.datasDisponiveis"] = bson.M{
				"$elemMatch": bson.M{"$regex": primitive.Regex{Pattern: monthPattern, Options: ""}},
			}
		}
	}
	if dto.Categoria != nil && *dto.Categoria != "" {
		filter["categoria"] = bson.M{"$regex": primitive.Regex{Pattern: string(*dto.Categoria), Options: "i"}}
	}

	cursor, err := db.Collection("destinos").Find(context.Background(), filter)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao buscar destinos.")
		return
	}
	defer cursor.Close(context.Background())

	var results []Destino
	if err = cursor.All(context.Background(), &results); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao decodificar destinos.")
		return
	}
	if results == nil {
		results = []Destino{}
	}
	RespondWithJSON(w, http.StatusOK, results)
}

func destinosPorCategoriaHandler(w http.ResponseWriter, r *http.Request) {

	type CategoriaCount struct {
		Categoria  Categoria `json:"categoria"`
		Quantidade int64     `json:"quantidade"`
	}
	results := []CategoriaCount{}

	// Using a WaitGroup if parallel execution is desired, but for a small list of categories, sequential is fine.
	var wg sync.WaitGroup
	var mu sync.Mutex // To protect shared 'results' slice if running in parallel

	for _, cat := range AllCategorias {
		wg.Add(1)
		go func(c Categoria) {
			defer wg.Done()
			count, err :=  db.Collection("destinos").CountDocuments(context.Background(), bson.M{"categoria": c})
			if err != nil {
				log.Printf("Erro ao contar documentos para categoria %s: %v", c, err)
				// Optionally handle error, e.g., return count 0 or skip
				return
			}
			mu.Lock()
			results = append(results, CategoriaCount{Categoria: c, Quantidade: count})
			mu.Unlock()
		}(cat)
	}
	wg.Wait()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
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

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"erro": "Erro interno ao gerar resposta JSON"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}


func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"erro": message})
}