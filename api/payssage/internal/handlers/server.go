// Package handlers implements the generated StrictServerInterface by
// aggregating one subpackage per feature. Each feature package owns the
// methods of its feature; this package only wires them together. Auth,
// validation, and error mapping run in the strict middleware stack (see
// internal/app); the handlers here are pure domain logic + fun envelope
// construction.
package handlers

import (
	"payssage/internal/handlers/collectors"
	"payssage/internal/handlers/intents"
	"payssage/internal/handlers/oauth"
	"payssage/internal/handlers/orgs"
	"payssage/internal/handlers/sellers"
	"payssage/internal/handlers/wallets"
	"payssage/internal/handlers/webhook_deliveries"
	"payssage/internal/handlers/webhook_endpoints"
	"payssage/internal/handlers/webhook_events"
	"payssage/internal/handlers/webhooks"
	"payssage/internal/services"
)

// Server implements openapi.StrictServerInterface by embedding every
// feature's Handlers: method promotion satisfies the generated interface
// with no delegation glue. The aliases exist only to give each feature's
// Handlers type a unique embeddable field name — embedding two
// *X.Handlers types directly would collide on the field name "Handlers".
type (
	OrgHandlers             = orgs.Handlers
	WalletHandlers          = wallets.Handlers
	OAuthHandlers           = oauth.Handlers
	CollectorHandlers       = collectors.Handlers
	SellerHandlers          = sellers.Handlers
	IntentHandlers          = intents.Handlers
	WebhookHandlers         = webhooks.Handlers
	WebhookEndpointHandlers = webhook_endpoints.Handlers
	WebhookEventHandlers    = webhook_events.Handlers
	WebhookDeliveryHandlers = webhook_deliveries.Handlers
)

type Server struct {
	*OrgHandlers
	*WalletHandlers
	*OAuthHandlers
	*CollectorHandlers
	*SellerHandlers
	*IntentHandlers
	*WebhookHandlers
	*WebhookEndpointHandlers
	*WebhookEventHandlers
	*WebhookDeliveryHandlers
}

// NewServer wires the per-feature handlers from the services aggregate.
func NewServer(ops *services.Operations) *Server {
	return &Server{
		OrgHandlers:             orgs.New(ops.Organizations),
		WalletHandlers:          wallets.New(ops.Wallets),
		OAuthHandlers:           oauth.New(ops.OAuth),
		CollectorHandlers:       collectors.New(ops.Collectors),
		SellerHandlers:          sellers.New(ops.Sellers),
		IntentHandlers:          intents.New(ops.Intents),
		WebhookHandlers:         webhooks.New(ops.Webhooks),
		WebhookEndpointHandlers: webhook_endpoints.New(ops.WebhookEndpoints),
		WebhookEventHandlers:    webhook_events.New(ops.WebhookEvents),
		WebhookDeliveryHandlers: webhook_deliveries.New(ops.WebhookDeliveries),
	}
}
