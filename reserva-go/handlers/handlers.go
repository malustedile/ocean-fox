package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reserva-go/services"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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
	Categoria *string `json:"categoria"`
}


func HelloHandler(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, "Hello Go Server!")
}

func MinhasReservasHandler(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("sessionId")
	if err != nil {
		if err == http.ErrNoCookie {
			RespondWithError(w, http.StatusUnauthorized, "Cookie 'sessionId' não encontrado.")
			return
		}
		RespondWithError(w, http.StatusBadRequest, "Erro ao ler cookie.")
		return
	}
	sessionID := sessionCookie.Value

	filter := bson.M{"sessionId": sessionID}
	cursor, err := services.ReservasCollection.Find(context.Background(), filter)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao buscar reservas.")
		return
	}
	defer cursor.Close(context.Background())

	var results []ReservaDocument
	if err = cursor.All(context.Background(), &results); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao decodificar reservas.")
		return
	}
	if results == nil {
		results = []ReservaDocument{} // Return empty array instead of null
	}
	RespondWithJSON(w, http.StatusOK, results)
}

func CriarDestinoHandler(w http.ResponseWriter, r *http.Request) {
	var dto Destino
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido.")
		return
	}

	if dto.Nome == "" || dto.Categoria == "" ||
		dto.Descricao.Navio == "" || dto.Descricao.Embarque == "" || dto.Descricao.Desembarque == "" ||
		len(dto.Descricao.DatasDisponiveis) == 0 || len(dto.Descricao.LugaresVisitados) == 0 ||
		dto.Descricao.Noites <= 0 || dto.Descricao.ValorPorPessoa <= 0 {
		RespondWithError(w, http.StatusBadRequest, "Campos obrigatórios ausentes ou inválidos.")
		return
	}
	dto.ID = primitive.NewObjectID() // Generate new ID for the document

	result, err := services.DestinosCollection.InsertOne(context.Background(), dto)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao adicionar destino.")
		return
	}
	RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"mensagem": "Destino adicionado com sucesso",
		"id":       result.InsertedID,
	})
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

	cursor, err := services.DestinosCollection.Find(context.Background(), filter)
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

func DestinosPorCategoriaHandler(w http.ResponseWriter, r *http.Request) {
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
			count, err := services.DestinosCollection.CountDocuments(context.Background(), bson.M{"categoria": c})
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

	RespondWithJSON(w, http.StatusOK, results)
}

func CancelarViagemHandler(w http.ResponseWriter, r *http.Request) {
	var dto CancelamentoDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido.")
		return
	}

	sessionCookie, err := r.Cookie("sessionId")
	if err != nil || sessionCookie.Value == "" {
		RespondWithError(w, http.StatusUnauthorized, "Cookie 'sessionId' inválido ou ausente.")
		return
	}
	sessionID := sessionCookie.Value

	if dto.ID == "" {
		RespondWithError(w, http.StatusBadRequest, "Campos obrigatórios ausentes ou inválidos.")
		return
	}

	log.Println("Cancelando reserva com ID:", dto.ID)
	objID, err := primitive.ObjectIDFromHex(dto.ID)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "ID de reserva inválido.")
		return
	}
	result := services.ReservasCollection.FindOne(context.Background(), bson.M{"_id": objID})
	if result == nil {
		RespondWithError(w, http.StatusNotFound, "Reserva não encontrada.")
	}
	var reserva ReservaDocument
	if err := result.Decode(&reserva); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao decodificar reserva.")
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
		RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
			"mensagem":      "Reserva cancelada com SUCESSO, mas FALHA ao notificar.",
		})
		return
	}

	err = services.RabbitMQChannelGlobal.PublishWithContext(
		context.Background(),
		services.ReservaCanceladaExchange, // exchange
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
		RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
			"mensagem":      "Reserva cancelada com SUCESSO, mas FALHA ao notificar (MQ).",
		})
		return
	}

	RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"mensagem":      "Reserva cancelada.",
	})

}

func ReservarDestinoHandler(w http.ResponseWriter, r *http.Request) {
	var dto ReservaDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Corpo da requisição inválido.")
		return
	}

	sessionCookie, err := r.Cookie("sessionId")
	if err != nil || sessionCookie.Value == "" {
		RespondWithError(w, http.StatusUnauthorized, "Cookie 'sessionId' inválido ou ausente.")
		return
	}
	sessionID := sessionCookie.Value

	if dto.Destino == "" || dto.DataEmbarque == "" || dto.NumeroPassageiros <= 0 ||
		dto.NumeroCabines <= 0 || dto.ValorTotal <= 0 {
		RespondWithError(w, http.StatusBadRequest, "Campos obrigatórios ausentes ou inválidos.")
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

	insertResult, err := services.ReservasCollection.InsertOne(context.Background(), reservaDoc)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Erro ao registrar reserva.")
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
		RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
			"mensagem":      "Reserva registrada com SUCESSO, mas FALHA ao notificar. Link de pagamento gerado.",
			"linkPagamento": linkPagamento,
			"reservaId":     insertResult.InsertedID,
		})
		return
	}

	err = services.RabbitMQChannelGlobal.PublishWithContext(
		context.Background(),
		services.ReservaCriadaExchange, // exchange
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
		RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
			"mensagem":      "Reserva registrada com SUCESSO, mas FALHA ao notificar (MQ). Link de pagamento gerado.",
			"linkPagamento": linkPagamento,
			"reservaId":     insertResult.InsertedID,
		})
		return
	}

	RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"mensagem":      "Reserva registrada. Link de pagamento gerado.",
		"linkPagamento": linkPagamento,
		"reservaId":     insertResult.InsertedID,
	})
}