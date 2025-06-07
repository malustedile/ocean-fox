package handlers

import (
	"encoding/json"
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

    RespondWithJSON(w, http.StatusOK, map[string]interface{}{
        "mensagem": "Reserva cancelada.",
    })
}