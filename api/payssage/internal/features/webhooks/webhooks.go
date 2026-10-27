// Package webhooks is the webhooks feature: the layers live in the
// subpackages, and Build wires the feature from here.
package webhooks

import (
	"net/http"

	"payssage/internal/features/webhooks/handlers"
	"payssage/internal/features/webhooks/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
	Handlers   = handlers.Handlers
)

func New(d Deps) *Operations { return service.New(d) }

func NewHandlers(ops *Operations) *Handlers { return handlers.New(ops) }

// RawRequestMiddleware captures the request body so the provider webhook
// receive can verify signatures against the exact bytes.
func RawRequestMiddleware(next http.Handler) http.Handler { return handlers.RawRequestMiddleware(next) }
