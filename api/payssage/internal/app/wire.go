package app

import (
	"net/http"
	"payssage/internal/authz"
	"payssage/internal/features/collectors"
	collectorsHandlers "payssage/internal/features/collectors/handlers"
	collectorsRepos "payssage/internal/features/collectors/repos"
	"payssage/internal/features/intents"
	intentsHandlers "payssage/internal/features/intents/handlers"
	intentsRepos "payssage/internal/features/intents/repos"
	"payssage/internal/features/oauth"
	oauthHandlers "payssage/internal/features/oauth/handlers"
	oauthRepos "payssage/internal/features/oauth/repos"
	"payssage/internal/features/orgs"
	orgsHandlers "payssage/internal/features/orgs/handlers"
	orgsRepos "payssage/internal/features/orgs/repos"
	"payssage/internal/features/providers"
	"payssage/internal/features/sellers"
	sellersHandlers "payssage/internal/features/sellers/handlers"
	sellersRepos "payssage/internal/features/sellers/repos"
	"payssage/internal/features/wallets"
	walletsHandlers "payssage/internal/features/wallets/handlers"
	walletsRepos "payssage/internal/features/wallets/repos"
	"payssage/internal/features/webhook_deliveries"
	webhookDeliveriesHandlers "payssage/internal/features/webhook_deliveries/handlers"
	webhookDeliveriesRepos "payssage/internal/features/webhook_deliveries/repos"
	"payssage/internal/features/webhook_endpoints"
	webhookEndpointsHandlers "payssage/internal/features/webhook_endpoints/handlers"
	webhookEndpointsRepos "payssage/internal/features/webhook_endpoints/repos"
	"payssage/internal/features/webhook_events"
	webhookEventsHandlers "payssage/internal/features/webhook_events/handlers"
	webhookEventsRepos "payssage/internal/features/webhook_events/repos"
	"payssage/internal/features/webhooks"
	webhooksHandlers "payssage/internal/features/webhooks/handlers"
	providers2 "payssage/internal/providers"
	"payssage/internal/sqlc"
	"payssage/ports"

	"github.com/jackc/pgx/v5"
	"github.com/riverqueue/river"
)

// ── Wire types ────────────────────────────────────────────────────────────

type repos struct {
	orgs       ports.OrganizationRepo
	wallets    ports.WalletRepo
	oauth      ports.OAuthStateRepo
	collectors ports.CollectorRepo
	sellers    ports.SellerRepo
	intents    ports.IntentRepo
	endpoints  ports.WebhookEndpointRepo
	deliveries ports.WebhookDeliveryRepo
	events     ports.WebhookEventRepo
}

type operations struct {
	orgs       *orgs.Operations
	wallets    *wallets.Operations
	oauth      *oauth.Operations
	collectors *collectors.Operations
	sellers    *sellers.Operations
	intents    *intents.Operations
	webhooks   *webhooks.Operations
	endpoints  *webhook_endpoints.Operations
	events     *webhook_events.Operations
	deliveries *webhook_deliveries.Operations
}

type middlewares struct {
	jwtAuth    func(http.Handler) http.Handler
	apiKeyAuth func(http.Handler) http.Handler
	anyAuth    func(http.Handler) http.Handler
}

type handlers struct {
	orgs       *orgsHandlers.Handlers
	wallets    *walletsHandlers.Handlers
	oauth      *oauthHandlers.Handlers
	collectors *collectorsHandlers.Handlers
	sellers    *sellersHandlers.Handlers
	intents    *intentsHandlers.Handlers
	webhooks   *webhooksHandlers.Handlers
	endpoints  *webhookEndpointsHandlers.Handlers
	events     *webhookEventsHandlers.Handlers
	deliveries *webhookDeliveriesHandlers.Handlers
}

// ── Init functions ────────────────────────────────────────────────────────

func (app *Payssage) initRepos(q *sqlc.Queries) repos {
	r := repos{
		orgs:       orgsRepos.NewRepo(q),
		wallets:    walletsRepos.NewRepo(q),
		oauth:      oauthRepos.NewRepo(q),
		collectors: collectorsRepos.NewRepo(q),
		sellers:    sellersRepos.NewRepo(q),
		intents:    intentsRepos.NewRepo(q),
		endpoints:  webhookEndpointsRepos.NewRepo(q),
		deliveries: webhookDeliveriesRepos.NewRepo(q),
		events:     webhookEventsRepos.NewRepo(q),
	}
	authz.Service = authz.New(r.orgs, r.wallets)
	return r
}

func (app *Payssage) initProviders(r repos) {
	mercadoPago := providers.NewMercadoPago(app.cfg.MercadoPagoConfig, r.intents, r.collectors, r.sellers, r.wallets, app.httpClient)

	providers2.PayssageProviders.OAuth = map[providers2.AvailableProviders]ports.OAuthProvider{
		providers2.MercadoPagoProvider: mercadoPago,
	}

	providers2.PayssageProviders.Payments = map[providers2.AvailableProviders]ports.PaymentAbstractionLayer{
		providers2.MercadoPagoProvider: mercadoPago,
	}

	providers2.PayssageProviders.Webhooks = map[providers2.AvailableProviders]ports.WebhookProvider{
		providers2.MercadoPagoProvider: mercadoPago,
	}
}

func (app *Payssage) initOperations(riverClient *river.Client[pgx.Tx], r repos) operations {
	return operations{
		orgs:       orgs.NewOperations(r.orgs, app.idxClient),
		wallets:    wallets.NewOperations(r.wallets, r.orgs),
		oauth:      oauth.NewOperations(r.wallets, r.orgs, r.oauth, r.collectors, r.sellers),
		collectors: collectors.NewOperations(r.collectors, r.orgs),
		sellers:    sellers.NewOperations(r.sellers, r.wallets, r.orgs),
		intents:    intents.NewOperations(r.intents, r.wallets, r.orgs, r.collectors, r.sellers),
		webhooks:   webhooks.NewOperations(riverClient, r.events, r.endpoints, r.deliveries),
		endpoints:  webhook_endpoints.NewOperations(r.endpoints, r.wallets, r.orgs),
		events:     webhook_events.NewOperations(r.events, r.wallets, r.orgs),
		deliveries: webhook_deliveries.NewOperations(r.deliveries, r.endpoints, r.wallets, r.orgs),
	}
}

func (app *Payssage) initMiddlewares() middlewares {
	var mw middlewares
	authMW := app.setupAuthMiddlewares()
	mw.jwtAuth = authMW.JWT()
	mw.apiKeyAuth = authMW.APIKey()
	mw.anyAuth = authMW.AnyAuth()
	return mw
}

func (app *Payssage) initHandlers(ops operations) handlers {
	return handlers{
		orgs:       orgsHandlers.NewHandlers(ops.orgs),
		wallets:    walletsHandlers.NewHandlers(ops.wallets),
		oauth:      oauthHandlers.NewHandlers(ops.oauth),
		collectors: collectorsHandlers.NewHandlers(ops.collectors),
		sellers:    sellersHandlers.NewHandlers(ops.sellers),
		intents:    intentsHandlers.NewHandlers(ops.intents),
		webhooks:   webhooksHandlers.NewHandlers(ops.webhooks),
		endpoints:  webhookEndpointsHandlers.NewHandlers(ops.endpoints),
		events:     webhookEventsHandlers.NewHandlers(ops.events),
		deliveries: webhookDeliveriesHandlers.NewHandlers(ops.deliveries),
	}
}
