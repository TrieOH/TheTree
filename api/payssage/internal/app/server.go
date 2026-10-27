package app

import (
	"payssage/internal/features/collectors"
	"payssage/internal/features/intents"
	"payssage/internal/features/oauth"
	"payssage/internal/features/orgs"
	"payssage/internal/features/sellers"
	"payssage/internal/features/wallets"
	"payssage/internal/features/webhook_deliveries"
	"payssage/internal/features/webhook_endpoints"
	"payssage/internal/features/webhook_events"
	"payssage/internal/features/webhooks"
	"payssage/internal/openapi"
)

// server implements openapi.StrictServerInterface by embedding every
// feature's Handlers: method promotion satisfies the generated interface
// with no delegation glue. The aliases only give each embedded field its
// own name — embedding two *X.Handlers directly would collide on
// "Handlers".
type (
	orgHandlers             = orgs.Handlers
	walletHandlers          = wallets.Handlers
	oauthHandlers           = oauth.Handlers
	collectorHandlers       = collectors.Handlers
	sellerHandlers          = sellers.Handlers
	intentHandlers          = intents.Handlers
	webhookHandlers         = webhooks.Handlers
	webhookEndpointHandlers = webhook_endpoints.Handlers
	webhookEventHandlers    = webhook_events.Handlers
	webhookDeliveryHandlers = webhook_deliveries.Handlers
)

type server struct {
	*orgHandlers
	*walletHandlers
	*oauthHandlers
	*collectorHandlers
	*sellerHandlers
	*intentHandlers
	*webhookHandlers
	*webhookEndpointHandlers
	*webhookEventHandlers
	*webhookDeliveryHandlers
}

var _ openapi.StrictServerInterface = (*server)(nil)
