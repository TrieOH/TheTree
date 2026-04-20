package app

import (
	"log"
	"net/http"
	"univents/internal/features/activities"
	"univents/internal/features/checkpoints"
	"univents/internal/features/editions"
	"univents/internal/features/events"
	"univents/internal/features/products"
	"univents/internal/features/purchases"
	"univents/internal/features/tickets"
	"univents/internal/interfaces/http/middleware"
	"univents/internal/interfaces/http/router"
	"univents/internal/interfaces/http/system"
	"univents/internal/platform/database"
	"univents/internal/platform/database/sqlc"
	"univents/internal/platform/queue"
	"univents/internal/platform/telemetry"
	"univents/internal/shared/ports"
	"univents/internal/shared/sockets"

	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type runtime struct {
	middlewares middlewares
	handlers    *router.HTTPDeps
	commands    commands
	queries     queries
	storeDeps   storeDeps
	repos       repos
	repoQueries *sqlc.Queries
	txRunner    database.TxRunner
	tracer      trace.Tracer
	logger      *zap.Logger
	asynq       asynqDeps
	ws          *sockets.Registry
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

type middlewares struct {
	authMW *middleware.AuthMiddleware
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
	rt.asynq = app.startAsynq(rt, rt.repos)
	rt.ws = sockets.New()
	defer app.stopAsynq(rt.asynq)
	rt.storeDeps = app.startStoreDeps()
	rt.commands = app.startCommands(rt, rt.repos)
	rt.queries = app.startQueries(rt, rt.repos)
	rt.handlers = app.startHandlers(rt)
	mux := router.CreateRouter(rt.handlers)
	port := viper.GetString("PORT")
	log.Printf("Univents listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
	return rt
}

func (app *Univents) startHandlers(rt runtime) *router.HTTPDeps {
	var handlers router.HTTPDeps
	handlers.AsynqmonHandler = asynqmon.New(asynqmon.Options{
		RootPath: "/admin/asynq",
		RedisConnOpt: asynq.RedisClientOpt{
			Addr:     viper.GetString("REDIS_ADDR"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
	})
	handlers.System = system.NewUniventsHandler()
	handlers.Events = events.NewHandler(rt.commands.events, rt.queries.events)
	handlers.Editions = editions.NewHandler(rt.commands.editions, rt.queries.editions)
	handlers.Activities = activities.NewHandler(rt.commands.activities, rt.queries.activities)
	handlers.Checkpoints = checkpoints.NewHandler(rt.commands.checkpoints, rt.queries.checkpoints)
	handlers.Tickets = tickets.NewHandler(rt.commands.tickets, rt.queries.tickets)
	handlers.Products = products.NewHandler(rt.commands.products, rt.queries.products)
	handlers.Purchases = purchases.NewHandler(rt.commands.purchases, rt.queries.purchases, rt.ws)

	handlers.AuthMiddleware = rt.middlewares.authMW
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
	cmd.products = products.NewCommandService(r.editions, r.products, r.purchases, app.payssage, rt.storeDeps.checkoutSessions, rt.ws, rt.storeDeps.publisher, app.minio, rt.asynq.client, rt.asynq.inspector, rt.tracer, app.sdbClient, rt.txRunner)
	cmd.purchases = purchases.NewCommandService(r.editions, r.products, r.purchases, app.payssage, rt.storeDeps.checkoutSessions, rt.ws, rt.storeDeps.publisher, app.minio, rt.asynq.client, rt.asynq.inspector, rt.tracer, app.sdbClient, rt.txRunner)

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

func (app *Univents) startMiddlewares(rt runtime) middlewares {
	var mw middlewares
	mw.authMW = middleware.NewAuthMiddleware(app.idxClient, rt.tracer)
	return mw
}

func (app *Univents) startAsynq(rt runtime, r repos) asynqDeps {
	var err error
	var deps asynqDeps
	workerHandlers := editions.NewAsynqService(r.editions, rt.tracer, rt.txRunner)
	activitiesAsyncHandlers := activities.NewAsynqService(r.activities, rt.tracer, rt.txRunner)
	purchasesAsyncHandlers := purchases.NewAsynqService(r.products, r.purchases, app.payssage, rt.storeDeps.publisher, rt.storeDeps.checkoutSessions, rt.ws, rt.tracer, rt.txRunner)
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
