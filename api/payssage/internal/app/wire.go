package app

import (
	"lib/database"
	"lib/errx"
	"lib/xslices"
	"net/http"
	"payssage/internal/database/sqlc"
	"payssage/internal/features/orgs"
	"payssage/ports"
	idx "sdk/identityx"
	"strings"

	mws "github.com/MintzyG/fun/middlewares"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ── Wire types ────────────────────────────────────────────────────────────

type repos struct {
	orgs ports.OrganizationRepo
	//intents             ports.IntentRepository
	//workspaces          ports.WorkspaceRepo
	//endpoints           ports.WebhookEndpointRepo
	//deliveries          ports.WebhookDeliveryRepo
	//events              ports.WebhookEventRepo
	//oauthStates         ports.OAuthStateRepo
	//providerCredentials ports.ProviderCredentialRepo
	//marketplaces        ports.MarketplaceConfigRepo
}

type queries struct {
	orgs *orgs.Queries
	//webhooks   *webhooks.QueryService
	//intents    *intents.QueryService
	//workspaces *workspaces.QueryService
	//apiKeys    *api_keys.QueryService
	//oauth      *oauth.QueryService
}

type commands struct {
	orgs *orgs.Commands
	//webhooks   *webhooks.CommandService
	//intents    *intents.CommandService
	//workspaces *workspaces.CommandService
	//apiKeys    *api_keys.CommandService
	//oauth      *oauth.CommandService
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
	orgs *orgs.Handlers
	//intents  *intents.Handler
	//wallets  *workspaces.Handler
	//webhooks *webhooks.Handler
	//oauth    *oauth.Handler
}

// ── Init functions ────────────────────────────────────────────────────────

func initRepos(q *sqlc.Queries, logger *zap.Logger, tracer trace.Tracer) repos {
	return repos{
		orgs: orgs.NewRepos(q, logger, tracer),
		//intents:             intents.NewIntentsRepo(q, logger, tracer),
		//workspaces:          workspaces.NewWorkspaceRepo(q, logger, tracer),
		//endpoints:           webhooks.NewWebhookEndpointRepo(q, logger, tracer),
		//deliveries:          webhooks.NewWebhookDeliveryRepo(q, logger, tracer),
		//events:              webhooks.NewWebhookEventRepo(q, logger, tracer),
		//oauthStates:         oauth.NewOAuthStatesRepo(q, logger, tracer),
		//providerCredentials: oauth.NewProviderCredentialsRepo(q, logger, tracer),
		//marketplaces:        oauth.NewMarketplaceConfigRepo(q, logger, tracer),
	}
}

func initQueries(r repos, idx *idx.Client, logger *zap.Logger, tx database.TxRunner, tracer trace.Tracer) queries {
	return queries{
		orgs: orgs.NewQueries(r.orgs, idx, logger, tracer, tx),
		//webhooks:   webhooks.NewQueryService(r.endpoints, r.deliveries, r.events, r.workspaces, logger, tx, tracer),
		//intents:    intents.NewQueryService(r.intents, r.workspaces, logger, tx, tracer),
		//workspaces: workspaces.NewQueryService(r.workspaces, logger, tx, tracer),
		//oauth:      oauth.NewQueryService(r.workspaces, r.marketplaces, logger, tx, tracer),
	}
}

func initCommands(r repos, idx *idx.Client, logger *zap.Logger, tx database.TxRunner, tracer trace.Tracer) commands {
	return commands{
		orgs: orgs.NewCommands(r.orgs, idx, logger, tracer, tx),
		//webhooks:   webhooks.NewCommandService(r.endpoints, r.deliveries, r.events, r.workspaces, r.intents, r.providerCredentials, river, logger, tx, tracer),
		//intents:    intents.NewCommandService(r.intents, r.workspaces, r.providerCredentials, r.marketplaces, cmd.webhooks, rt.paymentProviders.oauth, rt.paymentProviders.payments, logger, tx, tracer),
		//workspaces: workspaces.NewCommandService(r.workspaces, logger, tx, tracer),
		//oauth:      oauth.NewCommandService(r.intents, r.workspaces, r.oauthStates, r.providerCredentials, r.marketplaces, rt.paymentProviders.oauth, logger, tx, tracer),
	}
}

func initMiddlewares(logger *zap.Logger, cfg Config) middlewares {
	var mw middlewares
	authMW := setupAuthMiddlewares()
	mw.jwtAuth = authMW.JWT()
	mw.apiKeyAuth = authMW.APIKey()
	mw.anyAuth = authMW.AnyAuth()
	//mw.bodySize = mws.MaxBodySize(1 << 20)
	//mw.requestID = mws.RequestID(mws.RequestIDConfig{Header: "X-Request-ID"})
	mw.logger = mws.Logs(mws.Config{Logger: logger, SkipPrefixes: []string{"/metrics", "/health"}, RequestIDHeader: "X-Request-ID"})
	collectors, err := mws.NewCollectors(prometheus.DefaultRegisterer)
	if err != nil {
		errx.Exit(err, "Failed to create collectors")
	}
	mw.metrics = mws.Metrics(collectors, mws.MetricsConfig{SkipPrefixes: []string{"/metrics", "/health"}})
	mw.cors = mws.CORS(mws.CORSConfig{
		AllowedOrigins:   xslices.Clean(strings.Split(cfg.CorsAllowedOrigins, ",")),
		AllowedHeaders:   xslices.Clean(strings.Split(cfg.CorsAllowedHeaders, ",")),
		AllowCredentials: true,
	})
	//mw.realIP = mws.RealIP()
	//mw.recover = mws.Recover(logger)
	//mw.timeout = mws.Timeout(60 * time.Second)
	//mw.ratelimit = mws.RateLimit(mws.RateLimitConfig{RPS: 400, Burst: 20,
	//	KeyExtractor: func(r *http.Request) string { return r.RemoteAddr },
	//})
	return mw
}

func initHandlers(c commands, q queries) handlers {
	return handlers{
		orgs: orgs.NewHandlers(c.orgs, q.orgs),
		//intents:  intents.NewHandler(c.intents, q.intents),
		//wallets:  workspaces.NewHandler(c.workspaces, q.workspaces),
		//webhooks: webhooks.NewHandler(c.webhooks, q.webhooks),
		//oauth:    oauth.NewHandler(c.oauth, q.oauth),
	}
}
