package app

import (
	"Informd/internal/features/fields"
	"Informd/internal/features/forms"
	"Informd/internal/features/namespaces"
	"Informd/internal/features/responses"
	"Informd/internal/features/steps"
	"Informd/internal/openapi"
)

// server implements openapi.StrictServerInterface by embedding every
// feature's Handlers: method promotion satisfies the generated interface
// with no delegation glue. The aliases only give each embedded field its
// own name — embedding two *X.Handlers directly would collide on
// "Handlers".
type (
	namespaceHandlers = namespaces.Handlers
	formHandlers      = forms.Handlers
	stepHandlers      = steps.Handlers
	fieldHandlers     = fields.Handlers
	responseHandlers  = responses.Handlers
)

type server struct {
	*namespaceHandlers
	*formHandlers
	*stepHandlers
	*fieldHandlers
	*responseHandlers
}

var _ openapi.StrictServerInterface = (*server)(nil)
