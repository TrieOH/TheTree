package app

import (
	"lib/database"
	"payssage/internal/database/sqlc"
	"payssage/internal/features/api_keys"
	"payssage/internal/features/intents"
	"payssage/internal/features/oauth"
	"payssage/internal/features/webhooks"
	"payssage/internal/features/workspaces"
	"payssage/ports"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

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

type queries struct {
	webhooks   *webhooks.QueryService
	intents    *intents.QueryService
	workspaces *workspaces.QueryService
	apiKeys    *api_keys.QueryService
	oauth      *oauth.QueryService
}

type commands struct {
	webhooks   *webhooks.CommandService
	intents    *intents.CommandService
	workspaces *workspaces.CommandService
	apiKeys    *api_keys.CommandService
	oauth      *oauth.CommandService
}

type handlers struct {
	intents  *intents.Handler
	wallets  *workspaces.Handler
	webhooks *webhooks.Handler
	oauth    *oauth.Handler
}

// ── Init functions ────────────────────────────────────────────────────────

func initRepos(q *sqlc.Queries, logger *zap.Logger, tracer trace.Tracer) repos {
	return repos{
		intents:             intents.NewIntentsRepo(q, logger, tracer),
		workspaces:          workspaces.NewWorkspaceRepo(q, logger, tracer),
		apiKeys:             api_keys.NewApiKeyRepo(q, logger, tracer),
		endpoints:           webhooks.NewWebhookEndpointRepo(q, logger, tracer),
		deliveries:          webhooks.NewWebhookDeliveryRepo(q, logger, tracer),
		events:              webhooks.NewWebhookEventRepo(q, logger, tracer),
		oauthStates:         oauth.NewOAuthStatesRepo(q, logger, tracer),
		providerCredentials: oauth.NewProviderCredentialsRepo(q, logger, tracer),
		marketplaces:        oauth.NewMarketplaceConfigRepo(q, logger, tracer),
	}
}

func initQueries(r repos, logger *zap.Logger, tx database.TxRunner, tracer trace.Tracer) queries {
	return queries{
		webhooks:   webhooks.NewQueryService(r.endpoints, r.deliveries, r.events, r.workspaces, logger, tx, tracer),
		intents:    intents.NewQueryService(r.intents, r.workspaces, logger, tx, tracer),
		workspaces: workspaces.NewQueryService(r.workspaces, logger, tx, tracer),
		apiKeys:    api_keys.NewQueryService(r.apiKeys, r.workspaces, logger, tx, tracer),
		oauth:      oauth.NewQueryService(r.workspaces, r.marketplaces, logger, tx, tracer),
	}
}

func initCommands(r repos, river *river.Client[pgx.Tx], logger *zap.Logger, tx database.TxRunner, tracer trace.Tracer) commands {
	return commands{
		webhooks:   webhooks.NewCommandService(r.endpoints, r.deliveries, r.events, r.workspaces, r.intents, r.providerCredentials, river, logger, tx, tracer),
		intents:    intents.NewCommandService(r.intents, r.workspaces, r.providerCredentials, r.marketplaces, cmd.webhooks, rt.paymentProviders.oauth, rt.paymentProviders.payments, logger, tx, tracer),
		workspaces: workspaces.NewCommandService(r.workspaces, logger, tx, tracer),
		apiKeys:    api_keys.NewCommandService(r.apiKeys, r.workspaces, logger, tx, tracer),
		oauth:      oauth.NewCommandService(r.intents, r.workspaces, r.oauthStates, r.providerCredentials, r.marketplaces, rt.paymentProviders.oauth, logger, tx, tracer),
	}
}

func initHandlers(c commands, q queries) handlers {
	return handlers{
		intents:  intents.NewHandler(c.intents, q.intents),
		wallets:  workspaces.NewHandler(c.workspaces, q.workspaces),
		webhooks: webhooks.NewHandler(c.webhooks, q.webhooks),
		oauth:    oauth.NewHandler(c.oauth, q.oauth),
	}
}
