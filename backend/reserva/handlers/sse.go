package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
    sseClients   = make(map[string]chan string) 
    sseClientsMu sync.Mutex
)

// SSEHandler handles client connections for SSE
func SSEHandler(w http.ResponseWriter, r *http.Request) {

    fmt.Println(r.Cookie("sessionId"))
    cookie, err := r.Cookie("sessionId")
    if err != nil || cookie.Value == "" {
        http.Error(w, "Missing sessionId", http.StatusBadRequest)
        return
    }
    sessionId := cookie.Value

    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
        return
    }
	fmt.Println("New SSE client connected")
    // Set headers for SSE
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

    messageChan := make(chan string)
    sseClientsMu.Lock()
    sseClients[sessionId] = messageChan
    sseClientsMu.Unlock()

    // Remove client on disconnect
    defer func() {
        sseClientsMu.Lock()
        delete(sseClients, sessionId)
        sseClientsMu.Unlock()
        close(messageChan)
    }()

    // Send a ping every 30 seconds to keep connection alive
    pingDone := make(chan struct{})
    go func() {
        for {
            select {
            case <-pingDone:
                return
            case <-time.After(30 * time.Second):
                select {
                case messageChan <- ": ping\n":
                default:
                    return
                }
            }
        }
    }()

    // Listen for messages and send to client
    for {
        select {
        case msg := <-messageChan:
            fmt.Fprintf(w, "data: %s\n\n", msg)
            flusher.Flush()
        case <-r.Context().Done():
            close(pingDone) // sinaliza para a goroutine de ping parar
            return
        }
    }
}
// SSEMessage representa uma mensagem enviada para o cliente SSE
type SSEMessage struct {
    SessionID string          `json:"sessionId"`
    Msg       string          `json:"msg,omitempty"`
    EventType string          `json:"eventType"`
    Data      json.RawMessage `json:"data,omitempty"` // Campo opcional
}

// SSESendHandler envia uma mensagem para os clientes SSE conectados
func SSESendHandler(w http.ResponseWriter, r *http.Request) {
    var sseMsg SSEMessage
    if err := json.NewDecoder(r.Body).Decode(&sseMsg); err != nil {
        fmt.Println(err)
        http.Error(w, "Invalid JSON body", http.StatusBadRequest)
        return
    }
    if  sseMsg.EventType == "" || sseMsg.SessionID == "" {
        http.Error(w, "Missing required fields", http.StatusBadRequest)
        return
    }

    SendMessageToClient(sseMsg)

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Message sent"))
}

// SendMessageToClient envia a mensagem formatada para o cliente SSE
func SendMessageToClient(sseMsg SSEMessage) {
    sseClientsMu.Lock()
    defer sseClientsMu.Unlock()

    fmt.Println("Sending message to client:", sseMsg.SessionID)

    // Monta o payload com ou sem `data`
    var payloadMap = map[string]interface{}{
        "sessionId": sseMsg.SessionID,
        "msg":       sseMsg.Msg,
        "eventType": sseMsg.EventType,
    }
    fmt.Println(payloadMap)
    if len(sseMsg.Data) > 0 {
        var jsonData interface{}
        if err := json.Unmarshal(sseMsg.Data, &jsonData); err == nil {
            payloadMap["data"] = jsonData
        }
    }

    payloadBytes, err := json.Marshal(payloadMap)
    if err != nil {
        fmt.Println("Error marshalling payload:", err)
        return
    }

    if ch, ok := sseClients[sseMsg.SessionID]; ok {
        select {
        case ch <- string(payloadBytes):
            fmt.Println("Message sent to client:", sseMsg.SessionID)
        default:
            fmt.Println("Client channel full, skipping:", sseMsg.SessionID)
        }
    }
}