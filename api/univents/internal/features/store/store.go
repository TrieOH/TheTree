// Package store is the store feature: its handlers live here, its
// operations in service/, and Build wires the feature from here.
//
// The handlers serve the storefront's live-stock SSE surface (split 6):
// the raw `GET /editions/{edition_id}/store/stream` route (mounted by the
// app router, outside the strict handler — streaming must bypass the
// fun/validate envelope machinery).
package store

import (
	"univents/internal/features/store/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
