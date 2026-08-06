// Package handlers implements the generated StrictServerInterface by
// aggregating one subpackage per feature. Each feature package owns the
// methods of its feature; this package only wires them together. Auth,
// validation, and error mapping run in the strict middleware stack (see
// internal/app); the handlers here are pure domain logic + fun envelope
// construction.
package handlers

import (
	"univents/internal/handlers/badges"
	"univents/internal/handlers/certifications"
	"univents/internal/handlers/editions"
	"univents/internal/handlers/events"
	"univents/internal/handlers/products"
	"univents/internal/handlers/programs"
	"univents/internal/handlers/signatures"
	"univents/internal/handlers/ticket_types"
	"univents/internal/services"
)

// Server implements openapi.StrictServerInterface by embedding every
// feature's Handlers: method promotion satisfies the generated interface
// with no delegation glue. The aliases exist only to give each feature's
// Handlers type a unique embeddable field name — embedding two
// *X.Handlers types directly would collide on the field name "Handlers".
type (
	EventHandlers         = events.Handlers
	EditionHandlers       = editions.Handlers
	TicketTypeHandlers    = ticket_types.Handlers
	ProductHandlers       = products.Handlers
	ProgramHandlers       = programs.Handlers
	BadgeHandlers         = badges.Handlers
	SignatureHandlers     = signatures.Handlers
	CertificationHandlers = certifications.Handlers
)

type Server struct {
	*EventHandlers
	*EditionHandlers
	*TicketTypeHandlers
	*ProductHandlers
	*ProgramHandlers
	*BadgeHandlers
	*SignatureHandlers
	*CertificationHandlers
}

// NewServer wires the per-feature handlers from the services aggregate.
func NewServer(ops *services.Operations) *Server {
	return &Server{
		EventHandlers:         events.New(ops.Events),
		EditionHandlers:       editions.New(ops.Editions),
		TicketTypeHandlers:    ticket_types.New(ops.TicketTypes),
		ProductHandlers:       products.New(ops.Products),
		ProgramHandlers:       programs.New(ops.Programs),
		BadgeHandlers:         badges.New(ops.Badges),
		SignatureHandlers:     signatures.New(ops.Signatures),
		CertificationHandlers: certifications.New(ops.Certs),
	}
}
