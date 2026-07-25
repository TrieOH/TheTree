package app

import (
	"lib/database"
	"lib/errx"
	"lib/telemetry"
	"lib/xslices"
	"net/http"
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
	"strings"

	mws "github.com/MintzyG/fun/middlewares"
	"github.com/jackc/pgx/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/riverqueue/river"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
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
	logger            func(http.Handler) http.Handler
	cors              func(http.Handler) http.Handler
	jwtAuth           func(http.Handler) http.Handler
	apiKeyAuth        func(http.Handler) http.Handler
	anyAuth           func(http.Handler) http.Handler
	clientOnly        func(http.Handler) http.Handler
	projectClientOnly func(http.Handler) http.Handler
	metrics           func(http.Handler) http.Handler
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

func (app *Payssage) tracer() trace.Tracer {
	return otel.Tracer(app.cfg.AppName)
}

func (app *Payssage) txRunner() database.TxRunner {
	return database.NewPGXTxRunner(app.db)
}

func (app *Payssage) initRepos(q *sqlc.Queries) repos {
	return repos{
		orgs:       orgs.NewRepos(q, app.tracer()),
		wallets:    wallets.NewRepos(q, app.tracer()),
		oauth:      oauth.NewRepos(q, app.tracer()),
		collectors: collectors.NewRepos(q, app.tracer()),
		sellers:    sellers.NewRepos(q, app.tracer()),
		intents:    intents.NewRepos(q, app.tracer()),
		endpoints:  webhook_endpoints.NewRepos(q, app.tracer()),
		deliveries: webhook_deliveries.NewRepos(q, app.tracer()),
		events:     webhook_events.NewRepos(q, app.tracer()),
	}
}

func (app *Payssage) initProviders(r repos) {
	mercadoPago := providers.NewMercadoPago(app.cfg.MercadoPagoConfig, r.intents, r.collectors, r.sellers, r.wallets, app.tracer(), app.txRunner())

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
		orgs:       orgs.NewQueries(r.orgs, app.idxClient, app.tracer(), app.txRunner()),
		wallets:    wallets.NewQueries(r.wallets, r.orgs, app.tracer(), app.txRunner()),
		collectors: collectors.NewQueries(r.collectors, r.orgs, app.tracer(), app.txRunner()),
		sellers:    sellers.NewQueries(r.sellers, r.wallets, r.orgs, app.tracer(), app.txRunner()),
		intents:    intents.NewQueries(r.intents, r.wallets, r.orgs, app.tracer(), app.txRunner()),
		endpoints:  webhook_endpoints.NewQueries(r.endpoints, r.wallets, r.orgs, app.tracer(), app.txRunner()),
		events:     webhook_events.NewQueries(r.events, r.wallets, r.orgs, app.tracer(), app.txRunner()),
		deliveries: webhook_deliveries.NewQueries(r.deliveries, r.endpoints, r.wallets, r.orgs, app.tracer(), app.txRunner()),
	}
}

func (app *Payssage) initCommands(riverClient *river.Client[pgx.Tx], r repos) commands {
	return commands{
		orgs:      orgs.NewCommands(r.orgs, app.idxClient, app.tracer(), app.txRunner()),
		wallets:   wallets.NewCommands(r.wallets, r.orgs, app.tracer(), app.txRunner()),
		oauth:     oauth.NewCommands(r.wallets, r.orgs, r.oauth, r.collectors, r.sellers, app.tracer(), app.txRunner()),
		intents:   intents.NewCommands(r.intents, r.wallets, r.orgs, r.collectors, r.sellers, app.tracer(), app.txRunner()),
		webhooks:  webhooks.NewCommands(riverClient, r.events, r.endpoints, r.deliveries, app.tracer(), app.txRunner()),
		endpoints: webhook_endpoints.NewCommands(r.endpoints, r.wallets, r.orgs, app.tracer(), app.txRunner()),
	}
}

func (app *Payssage) initMiddlewares() middlewares {
	var mw middlewares
	authMW := app.setupAuthMiddlewares()
	mw.jwtAuth = authMW.JWT()
	mw.apiKeyAuth = authMW.APIKey()
	mw.anyAuth = authMW.AnyAuth()
	mw.logger = mws.Logs(mws.Config{Logger: telemetry.Log(), SkipPrefixes: []string{"/metrics", "/health"}, RequestIDHeader: "X-Request-ID"})
	collectors, err := mws.NewCollectors(prometheus.DefaultRegisterer)
	if err != nil {
		errx.Exit(err, "Failed to create collectors")
	}
	mw.metrics = mws.Metrics(collectors, mws.MetricsConfig{SkipPrefixes: []string{"/metrics", "/health"}})
	mw.cors = mws.CORS(mws.CORSConfig{
		AllowedOrigins:   xslices.Clean(strings.Split(app.cfg.CorsAllowedOrigins, ",")),
		AllowedHeaders:   xslices.Clean(strings.Split(app.cfg.CorsAllowedHeaders, ",")),
		AllowCredentials: true,
	})
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
