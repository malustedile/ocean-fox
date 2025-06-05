package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	reservaCriadaExchange   = "reserva-criada-exc"    // Fanout
	reservaCanceladaQueue   = "reserva-cancelada"     // Direct
	reservaCanceladaExchange = "reserva-cancelada-exc" // Direct
	pagamentoAprovadoExchange = "pagamento-aprovado-exc" // Direct
	pagamentoAprovadoRK     = "pagamento-aprovado"
	pagamentoRecusadoQueue  = "pagamento-recusado"
	pagamentoRecusadoExchange = "pagamento-recusado-exc" // Direct
	pagamentoRecusadoRK = "pagamento-recusado"
	bilheteGeradoQueue      = "bilhete-gerado"
)

// --- Structs (Data Models) ---

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

type DescricaoDestino struct {
	DatasDisponiveis []string `json:"datasDisponiveis" bson:"datasDisponiveis"`
	Navio            string   `json:"navio" bson:"navio"`
	Embarque         string   `json:"embarque" bson:"embarque"`
	Desembarque      string   `json:"desembarque" bson:"desembarque"`
	LugaresVisitados []string `json:"lugaresVisitados" bson:"lugaresVisitados"`
	Noites           int      `json:"noites" bson:"noites"`
	ValorPorPessoa   float64  `json:"valorPorPessoa" bson:"valorPorPessoa"`
}

type Destino struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Nome      string             `json:"nome" bson:"nome"`
	Categoria Categoria          `json:"categoria" bson:"categoria"`
	Descricao DescricaoDestino   `json:"descricao" bson:"descricao"`
}

type ReservaDTO struct { // For request body when creating a reserva
	Destino           string  `json:"destino"`
	DataEmbarque      string  `json:"dataEmbarque"`
	NumeroPassageiros int     `json:"numeroPassageiros"`
	NumeroCabines     int     `json:"numeroCabines"`
	ValorTotal        float64 `json:"valorTotal"`
}

type CancelamentoDTO struct { // For request body when canceling a reserva
	ID string `json:"id"` // Reserva ID to cancel
}

type ReservaDocument struct { // For MongoDB
	ID                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Destino           string             `json:"destino" bson:"destino"`
	SessionID         string             `json:"sessionId" bson:"sessionId"`
	DataEmbarque      string             `json:"dataEmbarque" bson:"dataEmbarque"`
	NumeroPassageiros int                `json:"numeroPassageiros" bson:"numeroPassageiros"`
	NumeroCabines     int                `json:"numeroCabines" bson:"numeroCabines"`
	ValorTotal        float64            `json:"valorTotal" bson:"valorTotal"`
	LinkPagamento     string             `json:"linkPagamento" bson:"linkPagamento"`
	Status            string             `json:"status" bson:"status"`
	PagamentoValido   *bool              `json:"pagamentoValido,omitempty" bson:"pagamentoValido,omitempty"` // Pointer to distinguish between false and not set
	CriadoEm          time.Time          `json:"criadoEm" bson:"criadoEm"`
}

type FiltrosDTO struct {
	Destino   *string    `json:"destino"`
	Mes       *string    `json:"mes"`
	Embarque  *string    `json:"embarque"`
	Categoria *Categoria `json:"categoria"`
}

// For RabbitMQ message content in pagamento-aprovado consumer
type PedidoReservaRMQ struct {
	ID string `json:"id"` // Assuming this is the string hex of ReservaDocument._id
	// ... other fields from reservaPayload if needed for signature, but JS example implies only 'id' is added
	// to the original 'reservaPayload' before being stringified and sent to 'reserva-criada-exc'.
	// The 'pagamento-aprovado' message structure seems to be { reserva: ..., assinatura: ... }
	// Let's define it more accurately based on the JS consumer.
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

// --- Global Variables (Collections, Channels, etc.) ---
var (
	mongoClient          *mongo.Client
	destinosCollection   *mongo.Collection
	reservasCollection   *mongo.Collection
	rabbitMQConnection   *amqp.Connection
	rabbitMQChannelGlobal *amqp.Channel // General purpose channel for publishing reserva-criada
)

// --- Helper Functions ---
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"erro": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
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

// --- Initialization Functions ---

func initMongoDB() {
	var err error
	clientOptions := options.Client().ApplyURI(mongoURI)
	mongoClient, err = mongo.Connect(context.Background(), clientOptions)
	failOnError(err, "Failed to connect to MongoDB")

	err = mongoClient.Ping(context.Background(), nil)
	failOnError(err, "Failed to ping MongoDB")
	log.Println("Connected to MongoDB!")

	db := mongoClient.Database(databaseName)
	destinosCollection = db.Collection(destinosCollectionName)
	reservasCollection = db.Collection(reservasCollectionName)
}

func initRabbitMQ() {
	var err error
	rabbitMQConnection, err = amqp.Dial(rabbitMQURI)
	failOnError(err, "Failed to connect to RabbitMQ")
	log.Println("Connected to RabbitMQ!")

	// General channel for publishing to "reserva-criada-exc"
	rabbitMQChannelGlobal, err = rabbitMQConnection.Channel()
	failOnError(err, "Failed to open a global RabbitMQ channel")
	
	// Declare queue "reserva-criada" (consumed by other services, durable)
	_, err = rabbitMQChannelGlobal.QueueDeclare(
		reservaCriadaQueue, // name
		true,               // durable
		false,              // delete when unused
		false,              // exclusive
		false,              // no-wait
		nil,                // arguments
	)
	failOnError(err, "Failed to declare queue 'reserva-criada'")

	// Declare exchange "reserva-criada-exc" (fanout, not durable as per JS)
	err = rabbitMQChannelGlobal.ExchangeDeclare(
		reservaCriadaExchange, // name
		"fanout",              // type
		false,                 // durable (JS has false)
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	failOnError(err, "Failed to declare exchange 'reserva-criada-exc'")

	// Channel for "reserva-cancelada"
	chReservaCancelada, err := rabbitMQConnection.Channel()
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
	err = rabbitMQChannelGlobal.ExchangeDeclare(
		reservaCanceladaExchange, // name
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
		reservaCanceladaExchange, // exchange
		false,                 // no-wait
		nil,                   // arguments
	)
	failOnError(err, "Failed to bind reserva cancelada queue to exchange")

	// Setup for Pagamento Aprovado consumer
	chPagamentoAprovado, err := rabbitMQConnection.Channel()
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
	chPagamentoRecusado, err := rabbitMQConnection.Channel()
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
	chBilheteGerado, err := rabbitMQConnection.Channel()
	failOnError(err, "Failed to open channel for bilhete gerado")
	defer chBilheteGerado.Close() // Close if only used for declaration here
	_, err = chBilheteGerado.QueueDeclare(bilheteGeradoQueue, true, false, false, false, nil)
	failOnError(err, "Failed to declare queue 'bilhete-gerado'")
}

// --- RabbitMQ Consumers ---

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
		reservasCollection.UpdateOne(
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
		reservasCollection.UpdateOne(
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
		err = rabbitMQChannelGlobal.PublishWithContext(
			context.Background(),
			reservaCanceladaExchange, // exchange
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

// --- HTTP Handlers ---

func helloHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, "Hello Go Server!")
}

func minhasReservasHandler(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("sessionId")
	if err != nil {
		if err == http.ErrNoCookie {
			respondWithError(w, http.StatusUnauthorized, "Cookie 'sessionId' não encontrado.")
			return
		}
		respondWithError(w, http.StatusBadRequest, "Erro ao ler cookie.")
		return
	}
	sessionID := sessionCookie.Value

	filter := bson.M{"sessionId": sessionID}
	cursor, err := reservasCollection.Find(context.Background(), filter)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erro ao buscar reservas.")
		return
	}
	defer cursor.Close(context.Background())

	var results []ReservaDocument
	if err = cursor.All(context.Background(), &results); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erro ao decodificar reservas.")
		return
	}
	if results == nil {
		results = []ReservaDocument{} // Return empty array instead of null
	}
	respondWithJSON(w, http.StatusOK, results)
}

func criarDestinoHandler(w http.ResponseWriter, r *http.Request) {
	var dto Destino
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido.")
		return
	}

	if dto.Nome == "" || dto.Categoria == "" ||
		dto.Descricao.Navio == "" || dto.Descricao.Embarque == "" || dto.Descricao.Desembarque == "" ||
		len(dto.Descricao.DatasDisponiveis) == 0 || len(dto.Descricao.LugaresVisitados) == 0 ||
		dto.Descricao.Noites <= 0 || dto.Descricao.ValorPorPessoa <= 0 {
		respondWithError(w, http.StatusBadRequest, "Campos obrigatórios ausentes ou inválidos.")
		return
	}
	dto.ID = primitive.NewObjectID() // Generate new ID for the document

	result, err := destinosCollection.InsertOne(context.Background(), dto)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erro ao adicionar destino.")
		return
	}
	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"mensagem": "Destino adicionado com sucesso",
		"id":       result.InsertedID,
	})
}

func buscarDestinosHandler(w http.ResponseWriter, r *http.Request) {
	var dto FiltrosDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido.")
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

	cursor, err := destinosCollection.Find(context.Background(), filter)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erro ao buscar destinos.")
		return
	}
	defer cursor.Close(context.Background())

	var results []Destino
	if err = cursor.All(context.Background(), &results); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erro ao decodificar destinos.")
		return
	}
	if results == nil {
		results = []Destino{}
	}
	respondWithJSON(w, http.StatusOK, results)
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
			count, err := destinosCollection.CountDocuments(context.Background(), bson.M{"categoria": c})
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

	respondWithJSON(w, http.StatusOK, results)
}

func cancelarViagemHandler(w http.ResponseWriter, r *http.Request) {
	var dto CancelamentoDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido.")
		return
	}

	sessionCookie, err := r.Cookie("sessionId")
	if err != nil || sessionCookie.Value == "" {
		respondWithError(w, http.StatusUnauthorized, "Cookie 'sessionId' inválido ou ausente.")
		return
	}
	sessionID := sessionCookie.Value

	if dto.ID == "" {
		respondWithError(w, http.StatusBadRequest, "Campos obrigatórios ausentes ou inválidos.")
		return
	}

	log.Println("Cancelando reserva com ID:", dto.ID)
	objID, err := primitive.ObjectIDFromHex(dto.ID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID de reserva inválido.")
		return
	}
	result := reservasCollection.FindOne(context.Background(), bson.M{"_id": objID})
	if result == nil {
		respondWithError(w, http.StatusNotFound, "Reserva não encontrada.")
	}
	var reserva ReservaDocument
	if err := result.Decode(&reserva); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erro ao decodificar reserva.")
		return
	}

	payloadParaFila := bson.M{
		"id":                dto.ID,
		"destino":           reserva.Destino,
		"sessionId":         sessionID,
		"dataEmbarque":      reserva.DataEmbarque,
		"numeroPassageiros": reserva.NumeroPassageiros,
		"numeroCabines":     reserva.NumeroCabines,
		"valorTotal":        reserva.ValorTotal,
	}

	reservaMsgBytes, err := json.Marshal(payloadParaFila)
	if err != nil {
		log.Printf("Erro ao fazer marshal da mensagem da reserva para RabbitMQ: %v", err)
		// Respond to client, but log the MQ error. The reservation is in DB.
		// Consider a retry mechanism for MQ or a compensating transaction.
		respondWithJSON(w, http.StatusCreated, map[string]interface{}{
			"mensagem":      "Reserva cancelada com SUCESSO, mas FALHA ao notificar.",
		})
		return
	}

	err = rabbitMQChannelGlobal.PublishWithContext(
		context.Background(),
		reservaCanceladaExchange, // exchange
		"",                    // routing key (fanout ignores this)
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         reservaMsgBytes,
		})
	if err != nil {
		log.Printf("Erro ao publicar mensagem de reserva criada: %v", err)
		// Similar to above, reservation is in DB.
		respondWithJSON(w, http.StatusCreated, map[string]interface{}{
			"mensagem":      "Reserva cancelada com SUCESSO, mas FALHA ao notificar (MQ).",
		})
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"mensagem":      "Reserva cancelada.",
	})

}

func reservarDestinoHandler(w http.ResponseWriter, r *http.Request) {
	var dto ReservaDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		respondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido.")
		return
	}

	sessionCookie, err := r.Cookie("sessionId")
	if err != nil || sessionCookie.Value == "" {
		respondWithError(w, http.StatusUnauthorized, "Cookie 'sessionId' inválido ou ausente.")
		return
	}
	sessionID := sessionCookie.Value

	if dto.Destino == "" || dto.DataEmbarque == "" || dto.NumeroPassageiros <= 0 ||
		dto.NumeroCabines <= 0 || dto.ValorTotal <= 0 {
		respondWithError(w, http.StatusBadRequest, "Campos obrigatórios ausentes ou inválidos.")
		return
	}

	linkPagamento := fmt.Sprintf("https://pagamento.fake/checkout?token=%s", uuid.NewString())

	reservaDoc := ReservaDocument{
		ID:                primitive.NewObjectID(),
		Destino:           dto.Destino,
		SessionID:         sessionID,
		DataEmbarque:      dto.DataEmbarque,
		NumeroPassageiros: dto.NumeroPassageiros,
		NumeroCabines:     dto.NumeroCabines,
		ValorTotal:        dto.ValorTotal,
		LinkPagamento:     linkPagamento,
		Status:            "AGUARDANDO_PAGAMENTO",
		CriadoEm:          time.Now().UTC(),
	}

	insertResult, err := reservasCollection.InsertOne(context.Background(), reservaDoc)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Erro ao registrar reserva.")
		return
	}

	// Prepare payload for RabbitMQ (what was signed and sent)
	// The JS code published {id: reserva.insertedId, ...reservaPayload}
	// where reservaPayload was the initial set of data for the document.
	// Here, reservaDoc already contains the ID and all fields.
	
	// Let's assume the `reserva-criada-exc` expects the *full* reserva document.
	// However, the JS code example shows:
	// JSON.stringify({ id: reserva.insertedId, ...reservaPayload })
	// where reservaPayload is the one created *before* insertOne, but with generated link and status.
	// For consistency, let's create what was sent.
	type ReservaPublicada struct {
		ID                string    `json:"id"` // ObjectID as hex string
		Destino           string    `json:"destino"`
		SessionID         string    `json:"sessionId"`
		DataEmbarque      string    `json:"dataEmbarque"`
		NumeroPassageiros int       `json:"numeroPassageiros"`
		NumeroCabines     int       `json:"numeroCabines"`
		ValorTotal        float64   `json:"valorTotal"`
		LinkPagamento     string    `json:"linkPagamento"`
		Status            string    `json:"status"`
		Bilhete           *string   `json:"bilhete"` // null in JS
		CriadoEm          string    `json:"criadoEm"` // ISOString
	}
	
	payloadParaFila := ReservaPublicada{
		ID: reservaDoc.ID.Hex(),
		Destino: reservaDoc.Destino,
		SessionID: reservaDoc.SessionID,
		DataEmbarque: reservaDoc.DataEmbarque,
		NumeroPassageiros: reservaDoc.NumeroPassageiros,
		NumeroCabines: reservaDoc.NumeroCabines,
		ValorTotal: reservaDoc.ValorTotal,
		LinkPagamento: reservaDoc.LinkPagamento,
		Status: reservaDoc.Status,
		Bilhete: nil, // Explicitly nil as in JS
		CriadoEm: reservaDoc.CriadoEm.Format(time.RFC3339Nano),
	}


	reservaMsgBytes, err := json.Marshal(payloadParaFila)
	if err != nil {
		log.Printf("Erro ao fazer marshal da mensagem da reserva para RabbitMQ: %v", err)
		// Respond to client, but log the MQ error. The reservation is in DB.
		// Consider a retry mechanism for MQ or a compensating transaction.
		respondWithJSON(w, http.StatusCreated, map[string]interface{}{
			"mensagem":      "Reserva registrada com SUCESSO, mas FALHA ao notificar. Link de pagamento gerado.",
			"linkPagamento": linkPagamento,
			"reservaId":     insertResult.InsertedID,
		})
		return
	}

	err = rabbitMQChannelGlobal.PublishWithContext(
		context.Background(),
		reservaCriadaExchange, // exchange
		"",                    // routing key (fanout ignores this)
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         reservaMsgBytes,
		})
	if err != nil {
		log.Printf("Erro ao publicar mensagem de reserva criada: %v", err)
		// Similar to above, reservation is in DB.
		respondWithJSON(w, http.StatusCreated, map[string]interface{}{
			"mensagem":      "Reserva registrada com SUCESSO, mas FALHA ao notificar (MQ). Link de pagamento gerado.",
			"linkPagamento": linkPagamento,
			"reservaId":     insertResult.InsertedID,
		})
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"mensagem":      "Reserva registrada. Link de pagamento gerado.",
		"linkPagamento": linkPagamento,
		"reservaId":     insertResult.InsertedID,
	})
}

// --- Main Function ---
func main() {
	initMongoDB()
	initRabbitMQ()  // This also starts the pagamentoAprovado consumer in a goroutine

	defer mongoClient.Disconnect(context.Background())
	defer rabbitMQConnection.Close()
	if rabbitMQChannelGlobal != nil {
		defer rabbitMQChannelGlobal.Close()
	}


	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/", helloHandler).Methods(http.MethodGet)
	r.HandleFunc("/minhas-reservas", minhasReservasHandler).Methods(http.MethodGet)
	r.HandleFunc("/destinos", criarDestinoHandler).Methods(http.MethodPost)
	r.HandleFunc("/destinos/buscar", buscarDestinosHandler).Methods(http.MethodPost)
	r.HandleFunc("/destinos-por-categoria", destinosPorCategoriaHandler).Methods(http.MethodGet)
	r.HandleFunc("/destinos/reservar", reservarDestinoHandler).Methods(http.MethodPost)
	r.HandleFunc("/destinos/cancelar", cancelarViagemHandler).Methods(http.MethodPost)

	// CORS middleware
	// AllowedOrigins can be more specific in production
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173"}), // Or specify your frontend domain(s)
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}),
		handlers.AllowCredentials(), // Important if you use cookies/auth headers from frontend
	)

	log.Printf("🦊 Go server is running at %s", httpServerAddress)
	if err := http.ListenAndServe(httpServerAddress, corsHandler(r)); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}