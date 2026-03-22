package sockets

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

type WSMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload,omitempty"`
}

type Registry struct {
	mu        sync.RWMutex
	conns     map[string]*websocket.Conn // key: sessionID
	callbacks map[string]func(WSMessage) // key: sessionID
}

func New() *Registry {
	return &Registry{
		conns:     make(map[string]*websocket.Conn),
		callbacks: make(map[string]func(WSMessage)),
	}
}
func (r *Registry) RegisterCallback(sessionID string, fn func(WSMessage)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.callbacks[sessionID] = fn
}

func (r *Registry) Register(sessionID string, conn *websocket.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.conns[sessionID] = conn
}

func (r *Registry) Notify(sessionID string, msg WSMessage) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	conn, ok := r.conns[sessionID]
	if !ok {
		return fmt.Errorf("no active connection for session %s", sessionID)
	}
	if cb, ok := r.callbacks[sessionID]; ok {
		cb(msg)
		return nil
	}
	return conn.WriteJSON(msg)
}

func (r *Registry) Remove(sessionID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.conns, sessionID)
	delete(r.callbacks, sessionID)
}

func MakeUpgrader() websocket.Upgrader {
	allowedOrigins := splitAndCleanCSV(viper.GetString("CORS_ALLOWED_ORIGINS"))
	if allowedOrigins == nil {
		log.Fatal("No AllowedOrigins set in CORS_ALLOWED_ORIGINS")
	}

	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowed[o] = struct{}{}
	}

	return websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" {
				return false
			}
			_, ok := allowed[origin]
			return ok
		},
	}
}

func splitAndCleanCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))

	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			out = append(out, v)
		}
	}

	if len(out) == 0 {
		return nil
	}

	return out
}
