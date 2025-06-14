package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reserva-go/services"
	"time"

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
	StatusPagamento   string     		 `json:"statusPagamento" bson:"statusPagamento"`
	Status			  string             `json:"status" bson:"status"` // Status of the reserva
	PagamentoValido   *bool              `json:"pagamentoValido,omitempty" bson:"pagamentoValido,omitempty"` // Pointer to distinguish between false and not set
	CriadoEm          time.Time          `json:"criadoEm" bson:"criadoEm"`
}

type Inscricao struct {
    SessionId string    `bson:"sessionId" json:"sessionId"`
    CriadoEm  time.Time `bson:"criadoEm" json:"criadoEm"`
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

func DestinosPorCategoriaHandler(w http.ResponseWriter, r *http.Request) {
    resp, err := http.Get("http://localhost:8080/destinos-por-categoria")
    if err != nil {
        http.Error(w, "Erro ao fazer a requisição: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Define o Content-Type como application/json
    w.Header().Set("Content-Type", "application/json")

    // Define o status code da resposta
    w.WriteHeader(resp.StatusCode)

    // Copia o corpo da resposta para o ResponseWriter
    _, err = io.Copy(w, resp.Body)
    if err != nil {
        http.Error(w, "Erro ao ler a resposta: "+err.Error(), http.StatusInternalServerError)
        return
    }
}

func BuscarDestinosHandler(w http.ResponseWriter, r *http.Request) {
    // Lê o corpo da requisição original
    bodyBytes, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "Erro ao ler corpo da requisição: "+err.Error(), http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    // Cria a requisição POST para o endpoint com o mesmo corpo
    req, err := http.NewRequest("POST", "http://localhost:8080/destinos/buscar", bytes.NewBuffer(bodyBytes))
    if err != nil {
        http.Error(w, "Erro ao criar requisição: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Define o Content-Type para JSON
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        http.Error(w, "Erro ao fazer requisição: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Define o Content-Type da resposta para JSON
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(resp.StatusCode)

    // Copia a resposta para o ResponseWriter
    _, err = io.Copy(w, resp.Body)
    if err != nil {
        http.Error(w, "Erro ao ler resposta: "+err.Error(), http.StatusInternalServerError)
        return
    }
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
