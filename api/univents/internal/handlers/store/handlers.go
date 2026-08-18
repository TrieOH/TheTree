// Package store serves the storefront's live-stock SSE surface (split 6):
// the raw `GET /editions/{edition_id}/store/stream` route (mounted by the
// app router, outside the strict handler — streaming must bypass the
// fun/validate envelope machinery).
package store

import (
	"univents/internal/services"
)

const module = "Univents"

type Handlers struct {
	ops *services.Store
}

func New(ops *services.Store) *Handlers { return &Handlers{ops: ops} }
