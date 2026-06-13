package app

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"lib/authz"
	database2 "lib/database"
	"lib/xslices"
	"univents/internal/features/activities"
	"univents/internal/features/checkpoints"
	"univents/internal/features/editions"
	"univents/internal/features/events"
	"univents/internal/features/products"
	"univents/internal/features/purchases"
	"univents/internal/features/security"
	"univents/internal/features/tickets"
	"univents/internal/platform/database"
	"univents/internal/platform/database/sqlc"
	"univents/internal/platform/queue"
	"univents/internal/platform/telemetry"
	"univents/internal/shared/errx"
	"univents/internal/shared/ports"
	"univents/internal/shared/sockets"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/middlewares"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	idx "sdk/identityx"
)

type runtime struct {
	middlewares mws
	handlers    *HTTPDeps
	commands    commands
	queries     queries
	storeDeps   storeDeps
	repos       repos
	repoQueries *sqlc.Queries
	txRunner    database2.TxRunner
	tracer      trace.Tracer
	logger      *zap.Logger
	asynq       asynqDeps
	wsRegistry  *sockets.Registry
}

type storeDeps struct {
	publisher        ports.InventoryPublisher
	subscriber       ports.InventorySubscriber
	checkoutSessions ports.PurchaseSessionStore
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

type repos struct {
	events      ports.EventsRepository
	editions    ports.EditionsRepository
	activities  ports.ActivitiesRepository
	checkpoints ports.CheckpointsRepository
	tickets     ports.TicketsRepository
	products    ports.ProductsRepository
	purchases   ports.PurchaseRepository
}

type mws struct {
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

type asynqDeps struct {
	client    *asynq.Client
	inspector *asynq.Inspector
	scheduler *asynq.Scheduler
	server    *asynq.Server
}

func (app *Univents) run() runtime {
	var rt runtime
	rt.repoQueries = sqlc.New(app.db)
	rt.txRunner = database.NewPGXTxRunner(app.db)
	rt.tracer = otel.Tracer(string(telemetry.UniventsTracer))
	rt.logger = telemetry.Log()

	rt.repos = app.startRepos(rt)
	rt.middlewares = app.startMiddlewares(rt)

	rt.storeDeps = app.startStoreDeps()
	rt.wsRegistry = sockets.New()

	rt.asynq = app.startAsynq(rt, rt.repos)
	defer app.stopAsynq(rt.asynq)

	rt.commands = app.startCommands(rt, rt.repos)
	rt.queries = app.startQueries(rt, rt.repos)
	rt.handlers = app.startHandlers(rt)

	mux := CreateRouter(rt.handlers)
	port := viper.GetString("PORT")
	log.Printf("Univents listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
	return rt
}

func (app *Univents) startHandlers(rt runtime) *HTTPDeps {
	var handlers HTTPDeps
	handlers.AsynqmonHandler = asynqmon.New(asynqmon.Options{
		RootPath: "/admin/asynq",
		RedisConnOpt: asynq.RedisClientOpt{
			Addr:     viper.GetString("REDIS_ADDR"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
	})
	handlers.Security = security.NewHandler(rt.wsRegistry)
	handlers.Events = events.NewHandler(rt.commands.events, rt.queries.events)
	handlers.Editions = editions.NewHandler(rt.commands.editions, rt.queries.editions)
	handlers.Activities = activities.NewHandler(rt.commands.activities, rt.queries.activities)
	handlers.Checkpoints = checkpoints.NewHandler(rt.commands.checkpoints, rt.queries.checkpoints)
	handlers.Tickets = tickets.NewHandler(rt.commands.tickets, rt.queries.tickets)
	handlers.Products = products.NewHandler(rt.commands.products, rt.queries.products)
	handlers.Purchases = purchases.NewHandler(rt.commands.purchases, rt.queries.purchases, rt.wsRegistry)

	handlers.BodySize = rt.middlewares.bodySize
	handlers.RequestID = rt.middlewares.requestID
	handlers.Logger = rt.middlewares.logger
	handlers.Metrics = rt.middlewares.metrics
	handlers.CORS = rt.middlewares.cors
	handlers.RealIP = rt.middlewares.realIP
	handlers.Recover = rt.middlewares.recover
	handlers.Timeout = rt.middlewares.timeout
	handlers.RateLimit = rt.middlewares.ratelimit
	handlers.Jwt = rt.middlewares.jwt
	handlers.ApiKey = rt.middlewares.apiKey
	handlers.AnyAuth = rt.middlewares.anyAuth
	handlers.AppName = app.cfg.AppName
	return &handlers
}

func (app *Univents) startQueries(rt runtime, r repos) queries {
	var q queries
	q.events = events.NewQueryService(r.events, rt.tracer, app.sdbClient, rt.txRunner)
	q.editions = editions.NewQueryService(r.events, r.editions, rt.tracer, app.sdbClient, rt.txRunner)
	q.activities = activities.NewQueryService(r.activities, r.editions, rt.tracer, app.sdbClient, rt.txRunner)
	q.checkpoints = checkpoints.NewQueryService(r.checkpoints, r.editions, rt.tracer, app.sdbClient, rt.txRunner)
	q.tickets = tickets.NewQueryService(r.tickets, r.editions, rt.tracer, rt.txRunner)
	q.products = products.NewQueryService(r.products, r.purchases, r.editions, rt.storeDeps.subscriber, rt.tracer, app.sdbClient, rt.txRunner)
	q.purchases = purchases.NewQueryService(r.products, r.purchases, r.editions, rt.storeDeps.subscriber, rt.tracer, app.sdbClient, rt.txRunner)
	return q
}

func (app *Univents) startCommands(rt runtime, r repos) commands {
	var cmd commands
	cmd.events = events.NewCommandService(r.events, app.minio, rt.tracer, app.sdbClient, rt.txRunner)
	cmd.editions = editions.NewCommandService(r.events, r.editions, rt.asynq.client, rt.tracer, app.sdbClient, rt.txRunner)
	cmd.activities = activities.NewCommandService(r.activities, r.editions, rt.asynq.client, rt.tracer, app.sdbClient, rt.txRunner)
	cmd.checkpoints = checkpoints.NewCommandService(r.checkpoints, r.editions, rt.asynq.client, rt.tracer, app.sdbClient, rt.txRunner)
	cmd.tickets = tickets.NewCommandService(r.editions, r.tickets, rt.asynq.client, rt.tracer, rt.txRunner)
	cmd.products = products.NewCommandService(r.editions, r.products, r.purchases, app.payssage, rt.storeDeps.checkoutSessions, rt.wsRegistry, rt.storeDeps.publisher, app.minio, rt.asynq.client, rt.asynq.inspector, rt.tracer, app.sdbClient, rt.txRunner)
	cmd.purchases = purchases.NewCommandService(r.editions, r.products, r.purchases, app.payssage, rt.storeDeps.checkoutSessions, rt.wsRegistry, rt.storeDeps.publisher, app.minio, rt.asynq.client, rt.asynq.inspector, rt.tracer, app.sdbClient, rt.txRunner)

	return cmd
}

func (app *Univents) startStoreDeps() storeDeps {
	var sd storeDeps
	sd.publisher = products.NewRedisInventoryPublisher(app.redis)
	sd.subscriber = products.NewRedisInventorySubscriber(app.redis)
	sd.checkoutSessions = purchases.NewCheckoutSessionStore(app.redis)
	return sd
}

func (app *Univents) startRepos(rt runtime) repos {
	var r repos
	r.events = events.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.editions = editions.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.activities = activities.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.checkpoints = checkpoints.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.tickets = tickets.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.products = products.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	r.purchases = purchases.NewRepo(rt.repoQueries, rt.logger, rt.tracer)
	return r
}

func (app *Univents) startMiddlewares(rt runtime) mws {
	var mw mws

	keyFunc := func(ctx context.Context, tokenStr string) (*idx.AccessClaims, error) {
		return app.idxClient.Tokens.VerifyAccessToken(ctx, tokenStr)
	}

	jwtHook := func(ctx context.Context, claims *idx.AccessClaims) (context.Context, error) {
		return authz.WithSubject(ctx, &authz.UserSubject{
			ID:    claims.Sub.ID,
			Email: claims.Sub.Email,
		}), nil
	}

	apiKeyHook := func(ctx context.Context, rawKey string) (context.Context, error) {
		telemetry.Log().Info("user tried to use api key",
			zap.String("message", "this service does not provide an api"),
			zap.String("key", rawKey),
		)
		return ctx, fun.ErrForbidden("this service does not provide access to its api")
	}

	authMW := middlewares.New[*idx.AccessClaims](keyFunc, jwtHook, apiKeyHook)
	mw.jwt = authMW.JWT()
	mw.apiKey = authMW.APIKey()
	mw.anyAuth = authMW.AnyAuth()
	mw.bodySize = middlewares.MaxBodySize(1 << 20)
	mw.requestID = middlewares.RequestID(middlewares.RequestIDConfig{Header: "X-Request-ID"})
	mw.logger = middlewares.Logs(middlewares.Config{Logger: rt.logger, SkipPrefixes: []string{"/admin/asynq"}, RequestIDHeader: "X-Request-ID"})
	collectors, err := middlewares.NewCollectors(prometheus.DefaultRegisterer)
	if err != nil {
		errx.Must(err, "Failed to create collectors")
	}
	mw.metrics = middlewares.Metrics(collectors, middlewares.MetricsConfig{SkipPrefixes: []string{"/metrics", "/health"}})
	mw.cors = middlewares.CORS(middlewares.CORSConfig{
		AllowedOrigins:   xslices.Clean(strings.Split(app.cfg.CorsAllowedOrigins, ",")),
		AllowCredentials: true,
	})
	mw.realIP = middlewares.RealIP()
	mw.recover = middlewares.Recover(rt.logger)
	mw.timeout = middlewares.Timeout(60 * time.Second)
	mw.ratelimit = middlewares.RateLimit(middlewares.RateLimitConfig{RPS: 400, Burst: 20,
		KeyExtractor: func(r *http.Request) string { return r.RemoteAddr },
	})
	return mw
}

func (app *Univents) startAsynq(rt runtime, r repos) asynqDeps {
	var err error
	var deps asynqDeps
	workerHandlers := editions.NewAsynqService(r.editions, rt.tracer, rt.txRunner)
	activitiesAsyncHandlers := activities.NewAsynqService(r.activities, rt.tracer, rt.txRunner)
	purchasesAsyncHandlers := purchases.NewAsynqService(r.products, r.purchases, app.payssage, rt.storeDeps.publisher, rt.storeDeps.checkoutSessions, rt.wsRegistry, rt.tracer, rt.txRunner)
	ticketsAsyncHandlers := tickets.NewAsynqService(r.tickets, r.products, r.activities, r.checkpoints, rt.tracer, app.sdbClient, rt.txRunner)
	deps.server, deps.client, deps.scheduler, deps.inspector, err = queue.InitAsynq(queue.Deps{
		Handlers:         workerHandlers,
		ActivityHandlers: activitiesAsyncHandlers,
		PurchaseHandlers: purchasesAsyncHandlers,
		TicketsHandler:   ticketsAsyncHandlers,
	})
	if err != nil {
		telemetry.Log().Fatal("failed to init Asynq", zap.Error(err))
	}
	return deps
}

func (app *Univents) stopAsynq(deps asynqDeps) {
	if err := deps.inspector.Close(); err != nil {
		telemetry.Log().Error("error closing the asynq inspector", zap.Error(err))
	}
	deps.scheduler.Shutdown()
	deps.server.Shutdown()
	if err := deps.client.Close(); err != nil {
		telemetry.Log().Error("error closing the asynq client", zap.Error(err))
	}
}
