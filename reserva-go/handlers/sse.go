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

type SSEMessage struct {
    SessionID string `json:"sessionId"`
    Msg       string `json:"msg"`
    EventType string `json:"eventType"`
}

// SSESendHandler sends a message to all connected SSE clients
func SSESendHandler(w http.ResponseWriter, r *http.Request) {
    var sseMsg SSEMessage
    if err := json.NewDecoder(r.Body).Decode(&sseMsg); err != nil {
        http.Error(w, "Invalid JSON body", http.StatusBadRequest)
        return
    }
    if sseMsg.Msg == "" || sseMsg.EventType == "" || sseMsg.SessionID == "" {
        http.Error(w, "Missing required fields", http.StatusBadRequest)
        return
    }
    payload := fmt.Sprintf(`{"sessionId":"%s","msg":"%s","eventType":"%s"}`, sseMsg.SessionID, sseMsg.Msg, sseMsg.EventType)

    sseClientsMu.Lock()
    if ch, ok := sseClients[sseMsg.SessionID]; ok {
        select {
        case ch <- payload:
        default:
        }
    }
    sseClientsMu.Unlock()

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Message sent"))
}