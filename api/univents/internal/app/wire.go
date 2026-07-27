package app

import (
	"lib/errx"
	"lib/objectstorage"
	"lib/telemetry"
	"lib/xslices"
	"net/http"
	"strings"
	"time"
	"univents/internal/features/editions"
	"univents/internal/features/products"
	"univents/internal/features/programs"
	"univents/internal/sqlc"

	idx "sdk/identityx"
	"univents/internal/features/events"
	"univents/internal/features/ticket_types"
	"univents/ports"

	mws "github.com/MintzyG/fun/middlewares"
	"github.com/prometheus/client_golang/prometheus"
)

// ── Wire types ────────────────────────────────────────────────────────────

type repos struct {
	events      ports.EventRepo
	editions    ports.EditionRepo
	ticketTypes ports.TicketTypeRepo
	products    ports.ProductRepo
	programs    ports.ProgramRepo
	occurrences ports.ProgramOccurrenceRepo
}

type queries struct {
	events      *events.Queries
	editions    *editions.Queries
	ticketTypes *ticket_types.Queries
	products    *products.Queries
	programs    *programs.Queries
}

type commands struct {
	events      *events.Commands
	editions    *editions.Commands
	ticketTypes *ticket_types.Commands
	products    *products.Commands
	programs    *programs.Commands
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
	events      *events.Handlers
	editions    *editions.Handlers
	ticketTypes *ticket_types.Handlers
	products    *products.Handlers
	programs    *programs.Handlers
}

// ── Init functions ────────────────────────────────────────────────────────

func initRepos(q *sqlc.Queries) repos {
	programsRepo := programs.NewRepos(q)
	return repos{
		events:      events.NewRepos(q),
		editions:    editions.NewRepos(q),
		ticketTypes: ticket_types.NewRepos(q),
		products:    products.NewRepos(q),
		programs:    programsRepo,
		occurrences: programsRepo,
	}
}

func initQueries(r repos) queries {
	return queries{
		events:      events.NewQueries(r.events),
		editions:    editions.NewQueries(r.events, r.editions),
		ticketTypes: ticket_types.NewQueries(r.editions, r.ticketTypes),
		products:    products.NewQueries(r.editions, r.products),
		programs:    programs.NewQueries(r.programs, r.occurrences),
	}
}

func initCommands(r repos, obj *objectstorage.Client, idx *idx.Client) commands {
	return commands{
		events:      events.NewCommands(r.events, obj, idx),
		editions:    editions.NewCommands(r.events, r.editions),
		ticketTypes: ticket_types.NewCommands(r.events, r.editions, r.ticketTypes),
		products:    products.NewCommands(r.events, r.editions, r.products),
		programs:    programs.NewCommands(r.events, r.editions, r.programs, r.occurrences),
	}
}

func initHandlers(q queries, c commands) handlers {
	return handlers{
		events:      events.NewHandlers(c.events, q.events),
		editions:    editions.NewHandlers(c.editions, q.editions),
		ticketTypes: ticket_types.NewHandlers(c.ticketTypes, q.ticketTypes),
		products:    products.NewHandlers(c.products, q.products),
		programs:    programs.NewHandlers(c.programs, q.programs),
	}
}

func initMiddlewares() middlewares {
	var mw middlewares
	authMW := SetupAuthMiddlewares()

	mw.jwt = authMW.JWT()
	mw.apiKey = authMW.APIKey()
	mw.anyAuth = authMW.AnyAuth()
	mw.bodySize = mws.MaxBodySize(1 << 20)
	mw.requestID = mws.RequestID(mws.RequestIDConfig{Header: "X-Request-ID"})
	mw.logger = mws.Logs(mws.Config{Logger: telemetry.Log(), SkipPrefixes: []string{"/health", "/metrics", "/admin/asynq"}, RequestIDHeader: "X-Request-ID"})
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
	mw.realIP = mws.RealIP()
	mw.recover = mws.Recover(telemetry.Log())
	mw.timeout = mws.Timeout(60 * time.Second)
	mw.ratelimit = mws.RateLimit(mws.RateLimitConfig{RPS: 400, Burst: 20,
		KeyExtractor: func(r *http.Request) string { return r.RemoteAddr },
	})
	return mw
}
