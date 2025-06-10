package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"reserva-go/services"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

    if !validaReservaDTO(dto) {
        RespondWithError(w, http.StatusBadRequest, "Campos obrigatórios ausentes ou inválidos.")
        return
    }

    ctx := r.Context()
    reservaDoc := NovaReservaDocument(dto, sessionID, "")

    insertResult, err := services.ReservasCollection.InsertOne(ctx, reservaDoc)
    if err != nil {
        RespondWithError(w, http.StatusInternalServerError, "Erro ao registrar reserva.")
        return
    }

    oid, ok := insertResult.InsertedID.(primitive.ObjectID)
    if !ok {
        RespondWithError(w, http.StatusInternalServerError, "Erro ao converter ID da reserva.")
        return
    }
    linkPagamento := gerarLinkPagamento(dto.ValorTotal, oid.Hex(), sessionID)
    if linkPagamento == "" {
        RespondWithError(w, http.StatusInternalServerError, "Erro ao gerar link de pagamento.")
        return
    }
    reservaDoc.LinkPagamento = linkPagamento
    _, err = services.ReservasCollection.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": reservaDoc})
    if err != nil {
        log.Printf("Erro ao atualizar reserva com link de pagamento: %v", err)
        ResponderReservaComAviso(w, linkPagamento, insertResult.InsertedID, "SUCESSO, mas FALHA ao atualizar reserva.")
        return
    }

    go tentaInscreverUsuario(sessionCookie)

    payloadParaFila := NovaReservaPublicada(reservaDoc)
    reservaMsgBytes, err := json.Marshal(payloadParaFila)
    if err != nil {
        log.Printf("Erro ao serializar mensagem da reserva: %v", err)
        ResponderReservaComAviso(w, linkPagamento, insertResult.InsertedID, "SUCESSO, mas FALHA ao notificar.")
        return
    }

    err = services.RabbitMQChannelGlobal.PublishWithContext(
        ctx,
        "reserva-criada-exc",
        "",
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        reservaMsgBytes,
        })
    if err != nil {
        log.Printf("Erro ao publicar mensagem de reserva criada: %v", err)
        ResponderReservaComAviso(w, linkPagamento, insertResult.InsertedID, "SUCESSO, mas FALHA ao notificar (MQ).")
        return
    }

    RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
        "mensagem":      "Reserva registrada. Link de pagamento gerado.",
        "linkPagamento": linkPagamento,
        "reservaId":     insertResult.InsertedID,
    })
}

// Funções auxiliares sugeridas:
func validaReservaDTO(dto ReservaDTO) bool {
    return dto.Destino != "" && dto.DataEmbarque != "" && dto.NumeroPassageiros > 0 &&
        dto.NumeroCabines > 0 && dto.ValorTotal > 0
}

func NovaReservaDocument(dto ReservaDTO, sessionID, linkPagamento string) ReservaDocument {
    return ReservaDocument{
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
}

type ReservaPublicada struct {
    ID                string      `json:"id"`
    Destino           string      `json:"destino"`
    SessionID         string      `json:"sessionId"`
    DataEmbarque      string      `json:"dataEmbarque"`
    NumeroPassageiros int         `json:"numeroPassageiros"`
    NumeroCabines     int         `json:"numeroCabines"`
    ValorTotal        float64     `json:"valorTotal"`
    LinkPagamento     string      `json:"linkPagamento"`
    Status            string      `json:"status"`
    Bilhete           interface{} `json:"bilhete"`
    CriadoEm          string      `json:"criadoEm"`
}

func NovaReservaPublicada(doc ReservaDocument) ReservaPublicada {
    return ReservaPublicada{
        ID:                doc.ID.Hex(),
        Destino:           doc.Destino,
        SessionID:         doc.SessionID,
        DataEmbarque:      doc.DataEmbarque,
        NumeroPassageiros: doc.NumeroPassageiros,
        NumeroCabines:     doc.NumeroCabines,
        ValorTotal:        doc.ValorTotal,
        LinkPagamento:     doc.LinkPagamento,
        Status:            doc.Status,
        Bilhete:           nil,
        CriadoEm:          doc.CriadoEm.Format(time.RFC3339Nano),
    }
}

func tentaInscreverUsuario(sessionCookie *http.Cookie) {
    req, err := http.NewRequest("POST", "http://localhost:3004/inscrever", nil)
    if err != nil {
        log.Printf("Erro ao criar requisição de inscrição: %v", err)
        return
    }
    req.AddCookie(sessionCookie)
    client := &http.Client{Timeout: 5 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Erro ao fazer requisição de inscrição: %v", err)
        return
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        log.Printf("Falha ao inscrever: status %d", resp.StatusCode)
    }
}

func ResponderReservaComAviso(w http.ResponseWriter, linkPagamento string, reservaId interface{}, aviso string) {
    RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
        "mensagem":      "Reserva registrada com " + aviso + " Link de pagamento gerado.",
        "linkPagamento": linkPagamento,
        "reservaId":     reservaId,
    })
}

func gerarLinkPagamento(valorTotal float64, IDReserva string, sessionID string) string {
    sistemaExternoURL := "http://localhost:3001/gerar-link"
    
    requestData := map[string]interface{}{
        "idReserva":  IDReserva,
        "valorTotal": valorTotal,
    }
    
    jsonData, err := json.Marshal(requestData)
    if err != nil {
        log.Println("Erro ao serializar dados da solicitação:", err)
        return ""
    }

    req, err := http.NewRequest("POST", sistemaExternoURL, bytes.NewBuffer(jsonData))
    if err != nil {
        log.Println("Erro ao criar requisição:", err)
        return ""
    }
    
    req.Header.Set("Content-Type", "application/json")
    req.AddCookie(&http.Cookie{Name: "sessionId", Value: sessionID})

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Println("Erro ao enviar requisição:", err)
        return ""
    }
    defer resp.Body.Close()

    var response map[string]interface{}
    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Println("Erro ao ler resposta do sistema externo:", err)
        return ""
    }
    
    if err := json.Unmarshal(bodyBytes, &response); err != nil {
        log.Println("Erro ao decodificar resposta do sistema externo:", err)
        return ""
    }

    log.Printf("Resposta do sistema externo: %v", response)
    link, ok := response["linkPagamento"].(string)
    if !ok {
        log.Println("Campo 'linkPagamento' não encontrado ou não é string")
        return ""
    }

    return link
}