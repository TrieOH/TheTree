// Package handlers implements the generated StrictServerInterface by
// aggregating one subpackage per feature. Each feature package owns the
// methods of its feature; this package only wires them together. Auth,
// validation, and error mapping run in the strict middleware stack (see
// internal/app); the handlers here are pure domain logic + fun envelope
// construction.
package handlers

import (
	"Informd/internal/handlers/fields"
	"Informd/internal/handlers/forms"
	"Informd/internal/handlers/namespaces"
	"Informd/internal/handlers/responses"
	"Informd/internal/handlers/steps"
	"Informd/internal/services"
)

// Server implements openapi.StrictServerInterface by embedding every
// feature's Handlers: method promotion satisfies the generated interface
// with no delegation glue. The aliases exist only to give each feature's
// Handlers type a unique embeddable field name — embedding two
// *X.Handlers types directly would collide on the field name "Handlers".
type (
	NamespaceHandlers = namespaces.Handlers
	FormHandlers      = forms.Handlers
	StepHandlers      = steps.Handlers
	FieldHandlers     = fields.Handlers
	ResponseHandlers  = responses.Handlers
)

type Server struct {
	*NamespaceHandlers
	*FormHandlers
	*StepHandlers
	*FieldHandlers
	*ResponseHandlers
}

// NewServer wires the per-feature handlers from the services aggregate.
func NewServer(ops *services.Operations) *Server {
	return &Server{
		NamespaceHandlers: namespaces.New(ops.Namespaces),
		FormHandlers:      forms.New(ops.Forms),
		StepHandlers:      steps.New(ops.Steps),
		FieldHandlers:     fields.New(ops.Fields),
		ResponseHandlers:  responses.New(ops.Responses),
	}
}
