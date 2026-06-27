package app

import "payssage/ports"

// ── Wire types ────────────────────────────────────────────────────────────

type repos struct {
	intents             ports.IntentRepository
	workspaces          ports.WorkspaceRepo
	apiKeys             ports.ApiKeysRepo
	endpoints           ports.WebhookEndpointRepo
	deliveries          ports.WebhookDeliveryRepo
	events              ports.WebhookEventRepo
	oauthStates         ports.OAuthStateRepo
	providerCredentials ports.ProviderCredentialRepo
	marketplaces        ports.MarketplaceConfigRepo
}
