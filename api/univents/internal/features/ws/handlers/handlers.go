// Package ws serves the realtime surfaces (split 6): the generated
// getWsToken op (REST) and the raw `WS /ws?token=...` route (mounted by the
// app router, outside the strict handler — a WebSocket handshake cannot
// carry Authorization headers and streaming must bypass the envelope
// machinery).
package handlers

import (
	"univents/internal/features/ws/service"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }
