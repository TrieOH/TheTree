// Package services aggregates every feature's operations layer. Import this
// package instead of the per-feature subpackages:
//
//	ops := services.NewOperations(r, riverClient, idxClient)
package services

import (
	"payssage/internal/repos"
	"payssage/internal/services/collectors"
	"payssage/internal/services/intents"
	"payssage/internal/services/oauth"
	"payssage/internal/services/orgs"
	"payssage/internal/services/sellers"
	"payssage/internal/services/wallets"
	"payssage/internal/services/webhook_deliveries"
	"payssage/internal/services/webhook_endpoints"
	"payssage/internal/services/webhook_events"
	"payssage/internal/services/webhooks"
	idx "sdk/identityx"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

// Type and constructor aliases for each feature's operations package.
type (
	Organizations     = orgs.Operations
	Wallets           = wallets.Operations
	OAuth             = oauth.Operations
	Collectors        = collectors.Operations
	Sellers           = sellers.Operations
	Intents           = intents.Operations
	Webhooks          = webhooks.Operations
	WebhookEndpoints  = webhook_endpoints.Operations
	WebhookEvents     = webhook_events.Operations
	WebhookDeliveries = webhook_deliveries.Operations
)

var (
	NewOrganizations     = orgs.NewOperations
	NewWallets           = wallets.NewOperations
	NewOAuth             = oauth.NewOperations
	NewCollectors        = collectors.NewOperations
	NewSellers           = sellers.NewOperations
	NewIntents           = intents.NewOperations
	NewWebhooks          = webhooks.NewOperations
	NewWebhookEndpoints  = webhook_endpoints.NewOperations
	NewWebhookEvents     = webhook_events.NewOperations
	NewWebhookDeliveries = webhook_deliveries.NewOperations
)

// Operations is the aggregate of every feature's operations, constructed
// once at startup and consumed by the HTTP handlers.
type Operations struct {
	Organizations     *Organizations
	Wallets           *Wallets
	OAuth             *OAuth
	Collectors        *Collectors
	Sellers           *Sellers
	Intents           *Intents
	Webhooks          *Webhooks
	WebhookEndpoints  *WebhookEndpoints
	WebhookEvents     *WebhookEvents
	WebhookDeliveries *WebhookDeliveries
}

// NewOperations wires every feature's operations from the shared repos.
func NewOperations(r *repos.Repos, riverClient *river.Client[pgx.Tx], idxClient *idx.Client) *Operations {
	return &Operations{
		Organizations:     NewOrganizations(r.Organizations, idxClient),
		Wallets:           NewWallets(r.Wallets, r.Organizations),
		OAuth:             NewOAuth(r.Wallets, r.Organizations, r.OAuth, r.Collectors, r.Sellers),
		Collectors:        NewCollectors(r.Collectors, r.Organizations),
		Sellers:           NewSellers(r.Sellers, r.Wallets, r.Organizations),
		Intents:           NewIntents(r.Intents, r.Wallets, r.Organizations, r.Collectors, r.Sellers),
		Webhooks:          NewWebhooks(riverClient, r.WebhookEvents, r.WebhookEndpoints, r.WebhookDeliveries),
		WebhookEndpoints:  NewWebhookEndpoints(r.WebhookEndpoints, r.Wallets, r.Organizations),
		WebhookEvents:     NewWebhookEvents(r.WebhookEvents, r.Wallets, r.Organizations),
		WebhookDeliveries: NewWebhookDeliveries(r.WebhookDeliveries, r.WebhookEndpoints, r.Wallets, r.Organizations),
	}
}
