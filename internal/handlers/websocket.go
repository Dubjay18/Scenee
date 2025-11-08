package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/Dubjay18/scenee/internal/auth"
	"gorm.io/gorm"
)

// WebSocketHandler handles WebSocket connections for real-time updates
type WebSocketHandler struct {
	DB          *gorm.DB
	connections map[string]*WSConnection
	mu          sync.RWMutex
}

// WSConnection represents a WebSocket connection
type WSConnection struct {
	UserID    string
	Messages  chan []byte
	closeChan chan struct{}
}

func NewWebSocketHandler(db *gorm.DB) *WebSocketHandler {
	return &WebSocketHandler{
		DB:          db,
		connections: make(map[string]*WSConnection),
	}
}

// Routes mounts the websocket routes
func (h *WebSocketHandler) Routes(r chi.Router) {
	r.Get("/ws", h.handleWebSocket)
}

// handleWebSocket handles WebSocket connections
// Note: This is a basic implementation. For production, consider using gorilla/websocket
func (h *WebSocketHandler) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	uid := auth.UserID(r.Context())
	if uid == "" {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
		return
	}

	// TODO: Implement actual WebSocket upgrade using gorilla/websocket
	// For now, this is a placeholder that uses Server-Sent Events (SSE) as an alternative

	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a connection
	conn := &WSConnection{
		UserID:    uid,
		Messages:  make(chan []byte, 10),
		closeChan: make(chan struct{}),
	}

	// Register connection
	h.mu.Lock()
	h.connections[uid] = conn
	h.mu.Unlock()

	// Clean up on disconnect
	defer func() {
		h.mu.Lock()
		delete(h.connections, uid)
		close(conn.Messages)
		h.mu.Unlock()
	}()

	// Send initial connection message
	initialMsg := map[string]interface{}{
		"type":      "connected",
		"message":   "WebSocket connection established",
		"timestamp": time.Now().Unix(),
	}
	if msgBytes, err := json.Marshal(initialMsg); err == nil {
		_, _ = w.Write([]byte("data: "))
		_, _ = w.Write(msgBytes)
		_, _ = w.Write([]byte("\n\n"))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}

	// Keep connection alive and send messages
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-conn.closeChan:
			return
		case msg := <-conn.Messages:
			_, _ = w.Write([]byte("data: "))
			_, _ = w.Write(msg)
			_, _ = w.Write([]byte("\n\n"))
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-ticker.C:
			// Send heartbeat
			heartbeat := map[string]interface{}{
				"type":      "heartbeat",
				"timestamp": time.Now().Unix(),
			}
			if msgBytes, err := json.Marshal(heartbeat); err == nil {
				_, _ = w.Write([]byte("data: "))
				_, _ = w.Write(msgBytes)
				_, _ = w.Write([]byte("\n\n"))
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			}
		}
	}
}

// BroadcastToUser sends a message to a specific user
func (h *WebSocketHandler) BroadcastToUser(userID string, messageType string, data interface{}) {
	h.mu.RLock()
	conn, exists := h.connections[userID]
	h.mu.RUnlock()

	if !exists {
		return
	}

	msg := map[string]interface{}{
		"type":      messageType,
		"data":      data,
		"timestamp": time.Now().Unix(),
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	select {
	case conn.Messages <- msgBytes:
	default:
		log.Printf("Message channel full for user %s", userID)
	}
}

// NotifyNewNotification sends a notification event to a user
func (h *WebSocketHandler) NotifyNewNotification(userID uuid.UUID, notification interface{}) {
	h.BroadcastToUser(userID.String(), "notification", notification)
}

// NotifyNewLike sends a like event to a user
func (h *WebSocketHandler) NotifyNewLike(userID uuid.UUID, likeData interface{}) {
	h.BroadcastToUser(userID.String(), "like", likeData)
}

// NotifyNewFollow sends a follow event to a user
func (h *WebSocketHandler) NotifyNewFollow(userID uuid.UUID, followData interface{}) {
	h.BroadcastToUser(userID.String(), "follow", followData)
}

// Mount returns a function that adds the routes under the given router
func (h *WebSocketHandler) Mount() func(r chi.Router) {
	return func(r chi.Router) { h.Routes(r) }
}
