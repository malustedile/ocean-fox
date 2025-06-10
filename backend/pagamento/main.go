package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	amqp "github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReservaPayload struct {
    ID               string `json:"id"`
    Destino          string `json:"destino"`
    DataEmbarque     string `json:"dataEmbarque"`
    NumeroPassageiros int    `json:"numeroPassageiros"`
    NumeroCabines    int    `json:"numeroCabines"`
    ValorTotal       float64 `json:"valorTotal"`
    LinkPagamento    string `json:"linkPagamento"`
    Status           string `json:"status"`
    CriadoEm         string `json:"criadoEm"`
}

type SolicitacaoPagamentoRequest struct {
    IDReserva  string  `json:"idReserva"`
    ValorTotal float64 `json:"valorTotal"`
}

type SolicitacaoPagamentoResponse struct {
    LinkPagamento string `json:"linkPagamento"`
    Status        string `json:"statusPagamento"`
}

type NotificacaoPagamento struct {
    ID         string  `json:"id"`
    Status     string  `json:"status"` // "PAGAMENTO_APROVADO" ou "PAGAMENTO_RECUSADO"
    ValorTotal float64 `json:"valorTotal"`
    IDReserva  string  `json:"idReserva"`
    SessionID  string  `json:"sessionId"`
}

var channelPagamentoAprovado *amqp.Channel
var channelPagamentoRecusado *amqp.Channel

var reservasCollection *mongo.Collection

func main() {
    // MongoDB
    ctx := context.Background()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:exemplo123@localhost:27017"))
    if err != nil {
        log.Fatal(err)
    }
    defer client.Disconnect(ctx)
    reservasCollection = client.Database("ocean-fox").Collection("reservas")

    // RabbitMQ
    rabbit, err := amqp.Dial("amqp://localhost")
    if err != nil {
        log.Fatal(err)
    }
    defer rabbit.Close()

    channelPagamentoAprovado, err = rabbit.Channel()
    if err != nil {
        log.Fatal("Erro ao criar canal pagamento aprovado:", err)
    }
    
    channelPagamentoRecusado, err = rabbit.Channel()
    if err != nil {
        log.Fatal("Erro ao criar canal pagamento recusado:", err)
    }
    
    pagamentoAprovadoExchange := "pagamento-aprovado-exc"

    _, err = channelPagamentoAprovado.QueueDeclare("pagamento-aprovado", true, false, false, false, nil)
    if err != nil {
        log.Fatal("Erro ao declarar queue pagamento aprovado:", err)
    }
    
    err = channelPagamentoAprovado.ExchangeDeclare(pagamentoAprovadoExchange, "direct", true, false, false, false, nil)
    if err != nil {
        log.Fatal("Erro ao declarar exchange pagamento aprovado:", err)
    }
    
    _, err = channelPagamentoRecusado.QueueDeclare("pagamento-recusado", true, false, false, false, nil)
    if err != nil {
        log.Fatal("Erro ao declarar queue pagamento recusado:", err)
    }
    
    r := mux.NewRouter()
    r.HandleFunc("/gerar-link", SolicitarLinkPagamentoHandler).Methods("POST")
    r.HandleFunc("/webhook/pagamento", WebhookPagamentoHandler).Methods("POST")

    port := "3001"
    fmt.Printf("App is running at 0.0.0.0:%s\n", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}

func SolicitarLinkPagamentoHandler(w http.ResponseWriter, r *http.Request) {
    log.Printf("gerando link de pagamento")
    var solicitacao SolicitacaoPagamentoRequest
    if err := json.NewDecoder(r.Body).Decode(&solicitacao); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    sessionID, err := r.Cookie("sessionId")
    if err != nil || sessionID.Value == "" {
        http.Error(w, "Session-ID header is required", http.StatusBadRequest)
        return
    }

    // Chama sistema externo de pagamento
    linkPagamento, err := solicitarLinkSistemaExterno(solicitacao, sessionID.Value)
    if err != nil {
        log.Printf("Erro ao solicitar link de pagamento: %v", err)
        http.Error(w, "Erro interno", http.StatusInternalServerError)
        return
    }

    response := SolicitacaoPagamentoResponse{
        LinkPagamento: linkPagamento,
        Status:        "AGUARDANDO_PAGAMENTO",
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func WebhookPagamentoHandler(w http.ResponseWriter, r *http.Request) {
    var notificacao NotificacaoPagamento
    if err := json.NewDecoder(r.Body).Decode(&notificacao); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    log.Printf("Webhook recebido - Reserva: %s, Status: %s", notificacao.IDReserva, notificacao.Status)

    // Atualiza status no MongoDB
    ctx := context.Background()
    objID, err := primitive.ObjectIDFromHex(notificacao.IDReserva)
    if err != nil {
        log.Printf("Erro ao converter ID da reserva: %v", err)
        http.Error(w, "ID inválido", http.StatusBadRequest)
        return
    }
    update := bson.M{
        "$set": bson.M{
            "statusPagamento": notificacao.Status,
        },
    }

    _, err = reservasCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
    if err != nil {
        log.Printf("Erro ao atualizar reserva no MongoDB: %v", err)
        http.Error(w, "Erro interno", http.StatusInternalServerError)
        return
    }

    // Publica mensagem na fila apropriada
    reserva := ReservaPayload{
        ID:         notificacao.IDReserva,
        ValorTotal: notificacao.ValorTotal,
        Status:     notificacao.Status,
    }
    reservaJSON, err := json.Marshal(reserva)
    if err != nil {
        log.Printf("Erro ao serializar reserva: %v", err)
        http.Error(w, "Erro interno", http.StatusInternalServerError)
        return
    }

    if notificacao.Status == "PAGAMENTO_APROVADO" {
        err := channelPagamentoAprovado.Publish(
            "pagamento-aprovado-exc",
            "pagamento-aprovado",
            false, false,
            amqp.Publishing{
                ContentType: "application/json",
                Body:        reservaJSON,
            },
        )
        if err != nil {
            log.Printf("Erro ao publicar pagamento aprovado: %v", err)
            http.Error(w, "Erro interno", http.StatusInternalServerError)
            return
        }
    } else if notificacao.Status == "PAGAMENTO_RECUSADO" {
        err := channelPagamentoRecusado.Publish(
            "",
            "pagamento-recusado",
            false, false,
            amqp.Publishing{
                ContentType: "application/json",
                Body:        reservaJSON,
            },
        )    
        if err != nil {
            log.Printf("Erro ao publicar pagamento recusado: %v", err)
        }
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Notificação processada com sucesso")
}

func solicitarLinkSistemaExterno(solicitacao SolicitacaoPagamentoRequest, sessionID string) (string, error) {
    sistemaExternoURL := "http://localhost:8000/link-pagamento"
    
    requestData := map[string]interface{}{
        "idReserva":  solicitacao.IDReserva,
        "valorTotal": solicitacao.ValorTotal,
    }
    
    jsonData, err := json.Marshal(requestData)
    if err != nil {
        fmt.Println("Erro ao serializar dados da solicitação:", err)
        return "", err
    }

    req, err := http.NewRequest("POST", sistemaExternoURL, bytes.NewBuffer(jsonData))
    if err != nil {
        fmt.Println("Erro ao criar requisição:", err)
        return "", err
    }
    
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Session-ID", sessionID)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Erro ao enviar requisição:", err)
        return "", err
    }
    defer resp.Body.Close()

    var response map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return "", err
    }

    link, ok := response["link"].(string)
    if !ok {
        return "", fmt.Errorf("resposta inválida do sistema externo")
    }

    return link, nil
}