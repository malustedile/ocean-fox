package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reserva-go/services"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

    log.Printf("[CancelarViagem] Cancelando reserva com ID: %s", dto.ID)
    objID, err := primitive.ObjectIDFromHex(dto.ID)
    if err != nil {
        RespondWithError(w, http.StatusBadRequest, "ID de reserva inválido.")
        return
    }

    result := services.ReservasCollection.FindOne(r.Context(), bson.M{"_id": objID})
    var reserva ReservaDocument
    if err := result.Decode(&reserva); err != nil {
        RespondWithError(w, http.StatusNotFound, "Reserva não encontrada.")
        return
    }

    // Atualiza o status para "cancelado"
    update := bson.M{"$set": bson.M{"status": "cancelado"}}
    _, err = services.ReservasCollection.UpdateOne(r.Context(), bson.M{"_id": objID}, update)
    if err != nil {
        log.Printf("[CancelarViagem] Erro ao atualizar status da reserva: %v", err)
        RespondWithError(w, http.StatusInternalServerError, "Erro ao cancelar a reserva.")
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
        log.Printf("[CancelarViagem] Erro ao serializar mensagem para RabbitMQ: %v", err)
        RespondWithJSON(w, http.StatusOK, map[string]interface{}{
            "mensagem": "Reserva cancelada com SUCESSO, mas FALHA ao notificar.",
        })
        return
    }

    err = services.RabbitMQChannelGlobal.PublishWithContext(
        r.Context(),
        "reserva-cancelada-exc",
        "",
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:        reservaMsgBytes,
        })
    if err != nil {
        log.Printf("[CancelarViagem] Erro ao publicar mensagem no RabbitMQ: %v", err)
        RespondWithJSON(w, http.StatusOK, map[string]interface{}{
            "mensagem": "Reserva cancelada com SUCESSO, mas FALHA ao notificar (MQ).",
        })
        return
    }
    enviarSse(map[string]interface{}{
        "id": dto.ID,
        "canceled": true,
    }, sessionID)
    RespondWithJSON(w, http.StatusOK, map[string]interface{}{
        "mensagem": "Reserva cancelada.",
    })
}


func enviarSse(data map[string]interface{}, sessionID string) (string, error) {
    fmt.Println("Sending SSE to update status")

    sistemaExternoURL := "http://localhost:3000/sse/send"
    
    requestData := map[string]interface{}{
        "sessionId": sessionID,
        "data": data,
        "eventType": "UPDATE_PAYMENT_STATUS",
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
    req.AddCookie(&http.Cookie{Name: "sessionId", Value: sessionID})

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