// Package repos aggregates every feature's repository layer. Import this
// package instead of the per-feature subpackages:
//
//	r := repos.New(q)
package repos

import (
	"payssage/internal/sqlc"

	"payssage/internal/repos/collectors"
	"payssage/internal/repos/intents"
	"payssage/internal/repos/oauth"
	"payssage/internal/repos/orgs"
	"payssage/internal/repos/sellers"
	"payssage/internal/repos/wallets"
	"payssage/internal/repos/webhook_deliveries"
	"payssage/internal/repos/webhook_endpoints"
	"payssage/internal/repos/webhook_events"
)

// Type and constructor aliases for each feature's repo package.
type (
	Collectors        = collectors.Repo
	Intents           = intents.Repo
	OAuth             = oauth.Repo
	Organizations     = orgs.Repo
	Sellers           = sellers.Repo
	Wallets           = wallets.Repo
	WebhookDeliveries = webhook_deliveries.Repo
	WebhookEndpoints  = webhook_endpoints.Repo
	WebhookEvents     = webhook_events.Repo
)

var (
	NewCollectors        = collectors.NewRepo
	NewIntents           = intents.NewRepo
	NewOAuth             = oauth.NewRepo
	NewOrganizations     = orgs.NewRepo
	NewSellers           = sellers.NewRepo
	NewWallets           = wallets.NewRepo
	NewWebhookDeliveries = webhook_deliveries.NewRepo
	NewWebhookEndpoints  = webhook_endpoints.NewRepo
	NewWebhookEvents     = webhook_events.NewRepo
)

// Repos is the aggregate of every feature repo, constructed once at startup.
type Repos struct {
	Organizations     *Organizations
	Wallets           *Wallets
	OAuth             *OAuth
	Collectors        *Collectors
	Sellers           *Sellers
	Intents           *Intents
	WebhookEndpoints  *WebhookEndpoints
	WebhookDeliveries *WebhookDeliveries
	WebhookEvents     *WebhookEvents
}

// New constructs every feature repo from the shared query handle.
func New(q *sqlc.Queries) *Repos {
	return &Repos{
		Organizations:     NewOrganizations(q),
		Wallets:           NewWallets(q),
		OAuth:             NewOAuth(q),
		Collectors:        NewCollectors(q),
		Sellers:           NewSellers(q),
		Intents:           NewIntents(q),
		WebhookEndpoints:  NewWebhookEndpoints(q),
		WebhookDeliveries: NewWebhookDeliveries(q),
		WebhookEvents:     NewWebhookEvents(q),
	}
}
