package app

import (
	"net/http"
	"payssage/internal/authz"
	"payssage/internal/features/collectors"
	"payssage/internal/features/intents"
	"payssage/internal/features/oauth"
	"payssage/internal/features/orgs"
	"payssage/internal/features/providers"
	"payssage/internal/features/sellers"
	"payssage/internal/features/wallets"
	"payssage/internal/features/webhook_deliveries"
	"payssage/internal/features/webhook_endpoints"
	"payssage/internal/features/webhook_events"
	"payssage/internal/features/webhooks"
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

type queries struct {
	orgs       *orgs.Queries
	wallets    *wallets.Queries
	collectors *collectors.Queries
	sellers    *sellers.Queries
	intents    *intents.Queries
	endpoints  *webhook_endpoints.Queries
	events     *webhook_events.Queries
	deliveries *webhook_deliveries.Queries
}

type commands struct {
	orgs      *orgs.Commands
	wallets   *wallets.Commands
	oauth     *oauth.Commands
	intents   *intents.Commands
	webhooks  *webhooks.Commands
	endpoints *webhook_endpoints.Commands
}

type middlewares struct {
	jwtAuth    func(http.Handler) http.Handler
	apiKeyAuth func(http.Handler) http.Handler
	anyAuth    func(http.Handler) http.Handler
}

type handlers struct {
	orgs       *orgs.Handlers
	wallets    *wallets.Handlers
	oauth      *oauth.Handlers
	collectors *collectors.Handlers
	sellers    *sellers.Handlers
	intents    *intents.Handlers
	webhooks   *webhooks.Handlers
	endpoints  *webhook_endpoints.Handlers
	events     *webhook_events.Handlers
	deliveries *webhook_deliveries.Handlers
}

// ── Init functions ────────────────────────────────────────────────────────

func (app *Payssage) initRepos(q *sqlc.Queries) repos {
	r := repos{
		orgs:       orgs.NewRepos(q),
		wallets:    wallets.NewRepos(q),
		oauth:      oauth.NewRepos(q),
		collectors: collectors.NewRepos(q),
		sellers:    sellers.NewRepos(q),
		intents:    intents.NewRepos(q),
		endpoints:  webhook_endpoints.NewRepos(q),
		deliveries: webhook_deliveries.NewRepos(q),
		events:     webhook_events.NewRepos(q),
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

func (app *Payssage) initQueries(r repos) queries {
	return queries{
		orgs:       orgs.NewQueries(r.orgs, app.idxClient),
		wallets:    wallets.NewQueries(r.wallets, r.orgs),
		collectors: collectors.NewQueries(r.collectors, r.orgs),
		sellers:    sellers.NewQueries(r.sellers, r.wallets, r.orgs),
		intents:    intents.NewQueries(r.intents, r.wallets, r.orgs),
		endpoints:  webhook_endpoints.NewQueries(r.endpoints, r.wallets, r.orgs),
		events:     webhook_events.NewQueries(r.events, r.wallets, r.orgs),
		deliveries: webhook_deliveries.NewQueries(r.deliveries, r.endpoints, r.wallets, r.orgs),
	}
}

func (app *Payssage) initCommands(riverClient *river.Client[pgx.Tx], r repos) commands {
	return commands{
		orgs:      orgs.NewCommands(r.orgs, app.idxClient),
		wallets:   wallets.NewCommands(r.wallets, r.orgs),
		oauth:     oauth.NewCommands(r.wallets, r.orgs, r.oauth, r.collectors, r.sellers),
		intents:   intents.NewCommands(r.intents, r.wallets, r.orgs, r.collectors, r.sellers),
		webhooks:  webhooks.NewCommands(riverClient, r.events, r.endpoints, r.deliveries),
		endpoints: webhook_endpoints.NewCommands(r.endpoints, r.wallets, r.orgs),
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

func (app *Payssage) initHandlers(c commands, q queries) handlers {
	return handlers{
		orgs:       orgs.NewHandlers(c.orgs, q.orgs),
		wallets:    wallets.NewHandlers(c.wallets, q.wallets),
		oauth:      oauth.NewHandlers(c.oauth),
		collectors: collectors.NewHandlers(q.collectors),
		sellers:    sellers.NewHandlers(q.sellers),
		intents:    intents.NewHandlers(c.intents, q.intents),
		webhooks:   webhooks.NewHandlers(c.webhooks),
		endpoints:  webhook_endpoints.NewHandlers(c.endpoints, q.endpoints),
		events:     webhook_events.NewHandlers(q.events),
		deliveries: webhook_deliveries.NewHandlers(q.deliveries),
	}
}
