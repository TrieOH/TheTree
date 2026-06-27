package app

import (
	"lib/database"
	"lib/errx"
	"lib/xslices"
	"net/http"
	"strings"
	"time"
	"univents/internal/features/activities"
	"univents/internal/features/checkpoints"
	"univents/internal/features/editions"
	"univents/internal/features/events"
	"univents/internal/features/products"
	"univents/internal/features/purchases"
	"univents/internal/features/security"
	"univents/internal/features/tickets"
	"univents/internal/platform/database/sqlc"
	"univents/internal/shared/ports"

	mws "github.com/MintzyG/fun/middlewares"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ── Wire types ────────────────────────────────────────────────────────────

type repos struct {
	events      ports.EventsRepository
	editions    ports.EditionsRepository
	activities  ports.ActivitiesRepository
	checkpoints ports.CheckpointsRepository
	tickets     ports.TicketsRepository
	products    ports.ProductsRepository
	purchases   ports.PurchaseRepository
}

type queries struct {
	events      *events.QueryService
	editions    *editions.QueryService
	activities  *activities.QueryService
	checkpoints *checkpoints.QueryService
	tickets     *tickets.QueryService
	products    *products.QueryService
	purchases   *purchases.QueryService
}

type commands struct {
	events      *events.CommandService
	editions    *editions.CommandService
	activities  *activities.CommandService
	checkpoints *checkpoints.CommandService
	tickets     *tickets.CommandService
	products    *products.CommandService
	purchases   *purchases.CommandService
}

type middlewares struct {
	logger    func(http.Handler) http.Handler
	requestID func(http.Handler) http.Handler
	bodySize  func(http.Handler) http.Handler
	metrics   func(http.Handler) http.Handler
	cors      func(http.Handler) http.Handler
	realIP    func(http.Handler) http.Handler
	recover   func(http.Handler) http.Handler
	timeout   func(http.Handler) http.Handler
	ratelimit func(http.Handler) http.Handler
	jwt       func(http.Handler) http.Handler
	apiKey    func(http.Handler) http.Handler
	anyAuth   func(http.Handler) http.Handler
}

type handlers struct {
	Events      *events.Handler
	Editions    *editions.Handler
	Activities  *activities.Handler
	Checkpoints *checkpoints.Handler
	Tickets     *tickets.Handler
	Products    *products.Handler
	Purchases   *purchases.Handler
	Security    *security.Handler
}

// ── Init functions ────────────────────────────────────────────────────────

func initRepos(q *sqlc.Queries, loggr *zap.Logger, tracer trace.Tracer) repos {
	return repos{
		events:      events.NewRepo(q, loggr, tracer),
		editions:    editions.NewRepo(q, loggr, tracer),
		activities:  activities.NewRepo(q, loggr, tracer),
		checkpoints: checkpoints.NewRepo(q, loggr, tracer),
		tickets:     tickets.NewRepo(q, loggr, tracer),
		products:    products.NewRepo(q, loggr, tracer),
		purchases:   purchases.NewRepo(q, loggr, tracer),
	}
}

func initQueries(r repos, tx database.TxRunner, loggr *zap.Logger, tracer trace.Tracer) queries {
	return queries{
		events:      events.NewQueryService(r.events, loggr, tracer, tx),
		editions:    editions.NewQueryService(r.events, r.editions, loggr, tracer, tx),
		activities:  activities.NewQueryService(r.activities, r.editions, loggr, tracer, tx),
		checkpoints: checkpoints.NewQueryService(r.checkpoints, r.editions, loggr, tracer, tx),
		tickets:     tickets.NewQueryService(r.tickets, r.editions, loggr, tracer, tx),
		products:    products.NewQueryService(r.products, r.purchases, r.editions, loggr, tracer, tx),
		purchases:   purchases.NewQueryService(r.products, r.purchases, r.editions, loggr, tracer, tx),
	}
}

func initCommands(r repos, tx database.TxRunner, loggr *zap.Logger, tracer trace.Tracer) commands {
	return commands{
		events:      events.NewCommandService(r.events, loggr, tracer, tx),
		editions:    editions.NewCommandService(r.events, r.editions, loggr, tracer, tx),
		activities:  activities.NewCommandService(r.activities, r.editions, loggr, tracer, tx),
		checkpoints: checkpoints.NewCommandService(r.checkpoints, r.editions, loggr, tracer, tx),
		tickets:     tickets.NewCommandService(r.editions, r.tickets, loggr, tracer, tx),
		products:    products.NewCommandService(r.editions, r.products, r.purchases, loggr, tracer, tx),
		purchases:   purchases.NewCommandService(r.editions, r.products, r.purchases, loggr, tracer, tx),
	}
}

func initHandlers(q queries, c commands) handlers {
	return handlers{
		//Security: security.NewHandler(rt.wsRegistry)
		Events:      events.NewHandler(c.events, q.events),
		Editions:    editions.NewHandler(c.editions, q.editions),
		Activities:  activities.NewHandler(c.activities, q.activities),
		Checkpoints: checkpoints.NewHandler(c.checkpoints, q.checkpoints),
		Tickets:     tickets.NewHandler(c.tickets, q.tickets),
		Products:    products.NewHandler(c.products, q.products),
		Purchases:   purchases.NewHandler(c.purchases, q.purchases),
	}
}

func initMiddlewares(logger *zap.Logger) middlewares {
	var mw middlewares
	authMW := SetupAuthMiddlewares()

	mw.jwt = authMW.JWT()
	mw.apiKey = authMW.APIKey()
	mw.anyAuth = authMW.AnyAuth()
	mw.bodySize = mws.MaxBodySize(1 << 20)
	mw.requestID = mws.RequestID(mws.RequestIDConfig{Header: "X-Request-ID"})
	mw.logger = mws.Logs(mws.Config{Logger: logger, SkipPrefixes: []string{"/admin/asynq"}, RequestIDHeader: "X-Request-ID"})
	collectors, err := mws.NewCollectors(prometheus.DefaultRegisterer)
	if err != nil {
		errx.Exit(err, "Failed to create collectors")
	}
	mw.metrics = mws.Metrics(collectors, mws.MetricsConfig{SkipPrefixes: []string{"/metrics", "/health"}})
	mw.cors = mws.CORS(mws.CORSConfig{
		AllowedOrigins:   xslices.Clean(strings.Split(app.cfg.CorsAllowedOrigins, ",")),
		AllowCredentials: true,
	})
	mw.realIP = mws.RealIP()
	mw.recover = mws.Recover(logger)
	mw.timeout = mws.Timeout(60 * time.Second)
	mw.ratelimit = mws.RateLimit(mws.RateLimitConfig{RPS: 400, Burst: 20,
		KeyExtractor: func(r *http.Request) string { return r.RemoteAddr },
	})
	return mw
}
