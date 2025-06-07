package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
    sseClients   = make(map[chan string]bool)
    sseClientsMu sync.Mutex
)

// SSEHandler handles client connections for SSE
func SSEHandler(w http.ResponseWriter, r *http.Request) {
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
    sseClients[messageChan] = true
    sseClientsMu.Unlock()

    // Remove client on disconnect
    defer func() {
        sseClientsMu.Lock()
        delete(sseClients, messageChan)
        sseClientsMu.Unlock()
        close(messageChan)
    }()

    // Send a ping every 30 seconds to keep connection alive
    go func() {
        for {
            time.Sleep(30 * time.Second)
            select {
            case messageChan <- ": ping\n":
            default:
                return
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
    if sseMsg.Msg == "" || sseMsg.EventType == "" {
        http.Error(w, "Missing required fields", http.StatusBadRequest)
        return
    }

    payload := fmt.Sprintf(`{"sessionId":"%s","msg":"%s","eventType":"%s"}`, sseMsg.SessionID, sseMsg.Msg, sseMsg.EventType)

    sseClientsMu.Lock()
    for ch := range sseClients {
        select {
        case ch <- payload:
        default:
        }
    }
    sseClientsMu.Unlock()

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Message sent"))
}