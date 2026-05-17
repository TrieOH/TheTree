package sockets

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type WSMessage struct {
	Type    string `json:"type"`
	Payload any    `json:"payload,omitempty"`
}

type Registry struct {
	mu    sync.RWMutex
	conns map[string]*websocket.Conn // key: sessionID
}

func New() *Registry {
	return &Registry{conns: make(map[string]*websocket.Conn)}
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
	return conn.WriteJSON(msg)
}

func (r *Registry) Remove(sessionID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.conns, sessionID)
}
