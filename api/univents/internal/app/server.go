package app

import (
	"univents/internal/features/badges"
	"univents/internal/features/certifications"
	"univents/internal/features/checkouts"
	"univents/internal/features/editions"
	"univents/internal/features/events"
	"univents/internal/features/payments"
	"univents/internal/features/products"
	"univents/internal/features/programs"
	"univents/internal/features/purchases"
	"univents/internal/features/signatures"
	"univents/internal/features/store"
	"univents/internal/features/ticket_types"
	"univents/internal/features/webhooks"
	"univents/internal/features/ws"
	"univents/internal/openapi"
)

// server implements openapi.StrictServerInterface by embedding every
// feature's Handlers: method promotion satisfies the generated interface
// with no delegation glue. The aliases only give each embedded field its
// own name — embedding two *X.Handlers directly would collide on
// "Handlers". The raw routes (ServeWS, ServeStoreStream) are promoted the
// same way.
type (
	eventHandlers         = events.Handlers
	editionHandlers       = editions.Handlers
	ticketTypeHandlers    = ticket_types.Handlers
	productHandlers       = products.Handlers
	programHandlers       = programs.Handlers
	badgeHandlers         = badges.Handlers
	signatureHandlers     = signatures.Handlers
	certificationHandlers = certifications.Handlers
	paymentHandlers       = payments.Handlers
	webhookHandlers       = webhooks.Handlers
	checkoutHandlers      = checkouts.Handlers
	purchaseHandlers      = purchases.Handlers
	wsHandlers            = ws.Handlers
	storeHandlers         = store.Handlers
)

type server struct {
	*eventHandlers
	*editionHandlers
	*ticketTypeHandlers
	*productHandlers
	*programHandlers
	*badgeHandlers
	*signatureHandlers
	*certificationHandlers
	*paymentHandlers
	*webhookHandlers
	*checkoutHandlers
	*purchaseHandlers
	*wsHandlers
	*storeHandlers
}

var _ openapi.StrictServerInterface = (*server)(nil)
