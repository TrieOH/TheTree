// Package ws serves the realtime surfaces (split 6): the generated
// getWsToken op (REST) and the raw `WS /ws?token=...` route (mounted by the
// app router, outside the strict handler — a WebSocket handshake cannot
// carry Authorization headers and streaming must bypass the envelope
// machinery).
package ws

import (
	"univents/internal/services"
)

type Handlers struct {
	ops *services.WS
}

func New(ops *services.WS) *Handlers { return &Handlers{ops: ops} }

const module = "Univents"
