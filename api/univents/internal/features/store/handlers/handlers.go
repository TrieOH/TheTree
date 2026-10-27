// Package store serves the storefront's live-stock SSE surface (split 6):
// the raw `GET /editions/{edition_id}/store/stream` route (mounted by the
// app router, outside the strict handler — streaming must bypass the
// fun/validate envelope machinery).
package handlers

import (
	"univents/internal/features/store/service"
)

const module = "Univents"

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }
