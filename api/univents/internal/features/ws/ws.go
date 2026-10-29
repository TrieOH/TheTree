// Package ws is the ws feature: its handlers live here, its operations in
// service/, and Build wires the feature from here.
//
// The handlers serve the realtime surfaces (split 6): the generated
// getWsToken op (REST) and the raw `WS /ws?token=...` route (mounted by the
// app router, outside the strict handler — a WebSocket handshake cannot
// carry Authorization headers and streaming must bypass the envelope
// machinery).
package ws

import (
	"univents/internal/features/ws/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
