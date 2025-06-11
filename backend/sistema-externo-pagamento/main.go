package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type NotificacaoPagamento struct {
    ID         string  `json:"id"`
    Status     string  `json:"status"` // "PAGAMENTO_APROVADO" ou "PAGAMENTO_RECUSADO"
    ValorTotal float64 `json:"valorTotal"`
    IDReserva  string  `json:"idReserva"`
    SessionID  string  `json:"sessionId"`
}

type LinkPagamentoResponse struct {
    Link        string                 `json:"link"`
    Notificacao NotificacaoPagamento   `json:"notificacao"`
}

type PagamentoRequest struct {
    IDReserva  string  `json:"idReserva"`
    ValorTotal float64 `json:"valorTotal"`
}

type ProcessarPagamentoRequest struct {

    IDReserva  string  `json:"idReserva"`
    ValorTotal float64 `json:"valorTotal"`
}

func LinkPagamentoHandler(w http.ResponseWriter, r *http.Request) {
    sessionId, err := r.Cookie("sessionId")
    if err != nil {
        http.Error(w, "Session-ID header is required", http.StatusBadRequest)
        return
    }
    var sessionID = sessionId.Value

    var pagamentoReq PagamentoRequest
    if err := json.NewDecoder(r.Body).Decode(&pagamentoReq); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    notificacao := NotificacaoPagamento{
        ID:         uuid.New().String(),
        IDReserva:  pagamentoReq.IDReserva,
        ValorTotal: pagamentoReq.ValorTotal,
        SessionID:  sessionID,
        Status:     "AGUARDANDO_PAGAMENTO",
    }

    response := LinkPagamentoResponse{
        Link:        fmt.Sprintf("https://pagamento.example.com/pay?sessionId=%s&idReserva=%s", sessionID, notificacao.IDReserva),
        Notificacao: notificacao,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
    log.Printf("Link de pagamento gerado: %s", response.Link)
    log.Printf("Notificação de pagamento: %+v", notificacao)
}

func ProcessarPagamentoHandler(w http.ResponseWriter, r *http.Request) {
    sessionId, err := r.Cookie("sessionId")
    if err != nil {
        http.Error(w, "Session-ID header is required", http.StatusBadRequest)
        return
    }
    var sessionID = sessionId.Value
    var req ProcessarPagamentoRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    idReserva := req.IDReserva
    valorTotal := strconv.FormatFloat(req.ValorTotal, 'f', 2, 64)
    
    if sessionID == "" || idReserva == "" || valorTotal == "" {
        http.Error(w, "sessionId and idReserva are required", http.StatusBadRequest)
        return
    }

    // Simula processamento assíncrono
    go processarPagamento(sessionID, idReserva, valorTotal)

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Pagamento sendo processado...")
}

// Processa o pagamento e envia webhook
func processarPagamento(sessionID, idReserva string, valorTotal string) {
    // Simula tempo de processamento
    time.Sleep(2 * time.Second)

    // Simula aprovação/recusa aleatória
    status := "PAGAMENTO_RECUSADO"
    if randBool() {
        status = "PAGAMENTO_APROVADO"
    }

    valor, err := strconv.ParseFloat(valorTotal, 64)
    if err != nil {
        log.Printf("Erro ao converter valorTotal: %v", err)
        return
    }

    notificacao := NotificacaoPagamento{
        ID:         uuid.New().String(),
        Status:     status,
        IDReserva:  idReserva,
        ValorTotal: valor,
        SessionID:  sessionID,
    }

    // Envia webhook para o MS Pagamento
    enviarWebhook(notificacao)
}

func enviarWebhook(notificacao NotificacaoPagamento) {
    webhookURL := "http://localhost:3001/webhook/pagamento" // URL do MS Pagamento
    
    jsonData, err := json.Marshal(notificacao)
    if err != nil {
        log.Printf("Erro ao serializar notificação: %v", err)
        return
    }

    resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
    log.Printf("Enviando webhook para %s", webhookURL)
    if err != nil {
        log.Printf("Erro ao enviar webhook: %v", err)
        return
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
        log.Printf("Webhook enviado com sucesso para reserva %s - Status: %s", 
            notificacao.IDReserva, notificacao.Status)
    } else {
        log.Printf("Erro no webhook - Status: %d", resp.StatusCode)
    }
}

func main() {
    r := mux.NewRouter()
    
    r.HandleFunc("/link-pagamento", LinkPagamentoHandler).Methods(http.MethodPost)
    r.HandleFunc("/pagar", ProcessarPagamentoHandler).Methods(http.MethodPost)


	// CORS middleware
	// AllowedOrigins can be more specific in production
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173"}), // Or specify your frontend domain(s)
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "X-Requested-With"}),
		handlers.AllowCredentials(), // Important if you use cookies/auth headers from frontend
	)

    if err := http.ListenAndServe(":8000", corsHandler(r)); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}

func randBool() bool {
    rand.Seed(time.Now().UnixNano())
    return rand.Intn(2) == 1
}