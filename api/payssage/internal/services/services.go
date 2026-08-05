// Package services aggregates every feature's operations layer. Import this
// package instead of the per-feature subpackages:
//
//	ops := services.NewOperations(r, riverClient, idxClient)
package services

import (
	"payssage/internal/authz"
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
// NewOperations wires every feature's operations from the shared repos.
// Authorization arrives by injection through the same seam — no
// service-locator globals.
func NewOperations(r *repos.Repos, authzSvc *authz.Service, riverClient *river.Client[pgx.Tx], idxClient *idx.Client) *Operations {
	return &Operations{
		Organizations:     NewOrganizations(r.Organizations, idxClient, authzSvc),
		Wallets:           NewWallets(r.Wallets, r.Organizations, authzSvc),
		OAuth:             NewOAuth(r.Wallets, r.Organizations, r.OAuth, r.Collectors, r.Sellers, authzSvc),
		Collectors:        NewCollectors(r.Collectors, r.Organizations, authzSvc),
		Sellers:           NewSellers(r.Sellers, r.Wallets, r.Organizations, authzSvc),
		Intents:           NewIntents(r.Intents, r.Wallets, r.Organizations, r.Collectors, r.Sellers, authzSvc),
		Webhooks:          NewWebhooks(riverClient, r.WebhookEvents, r.WebhookEndpoints, r.WebhookDeliveries, authzSvc),
		WebhookEndpoints:  NewWebhookEndpoints(r.WebhookEndpoints, r.Wallets, r.Organizations, authzSvc),
		WebhookEvents:     NewWebhookEvents(r.WebhookEvents, r.Wallets, r.Organizations, authzSvc),
		WebhookDeliveries: NewWebhookDeliveries(r.WebhookDeliveries, r.WebhookEndpoints, r.Wallets, r.Organizations, authzSvc),
	}
}
