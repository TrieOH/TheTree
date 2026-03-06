package initialization

import (
	"context"
	"log"
	"net/http"
	"time"
	productsCommands "univents/internal/commerce/application/product/commands"
	productsQueries "univents/internal/commerce/application/product/queries"
	ticketsCommands "univents/internal/commerce/application/ticket/commands"
	ticketsQueries "univents/internal/commerce/application/ticket/queries"
	commerceInfra "univents/internal/commerce/infrastructure"
	productshttp "univents/internal/commerce/interfaces/http/products"
	tickethttp "univents/internal/commerce/interfaces/http/tickets"
	activityAsync "univents/internal/core/application/activity/async"
	activityCommands "univents/internal/core/application/activity/commands"
	activityQueries "univents/internal/core/application/activity/queries"
	checkpointCommands "univents/internal/core/application/checkpoint/commands"
	"univents/internal/core/application/edition/async"
	editionCommands "univents/internal/core/application/edition/commands"
	editionQueries "univents/internal/core/application/edition/queries"
	"univents/internal/core/application/event/commands"
	"univents/internal/core/application/event/queries"
	"univents/internal/core/infrastructure"
	eventhttp "univents/internal/core/interfaces/http"
	activityhttp "univents/internal/core/interfaces/http/activities"
	checkpointshttp "univents/internal/core/interfaces/http/checkpoints"
	editionhttp "univents/internal/core/interfaces/http/editions"
	"univents/internal/interfaces/http/middleware"
	"univents/internal/interfaces/http/router"
	"univents/internal/interfaces/http/system"
	"univents/internal/plataform/database"
	"univents/internal/plataform/database/sqlc"
	"univents/internal/plataform/telemetry"
	"univents/internal/worker"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/go-co-op/gocron/v2"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type UniventsApp struct {
	Port      string
	DB        *pgxpool.Pool
	Redis     *redis.Client
	Scheduler gocron.Scheduler
	GaClient  *goauth.Client

	Deps *router.HTTPDeps
}

func UniventsSetup() *UniventsApp {
	var app UniventsApp

	LoadEnv(&app)
	SetupGoAuth(&app)
	SetupFUN()
	if viper.GetString("ENV") != "test" {
		SetupDB(&app, "./internal/plataform/database/migrations")
	} else {
		log.Println("WE'RE TESTING")
		SetupDB(&app, "../internal/plataform/database/migrations")
	}
	app.Redis = SetupRedis(15 * time.Second)
	SetupCron(app.DB, &app)

	return &app
}

func UniventsStart(app *UniventsApp, skipMux bool) {
	ctx := context.Background()

	defer app.DB.Close()
	defer app.Redis.Close()

	shutdown, err := telemetry.InitTracer(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func(ctx context.Context) {
		err := shutdown(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(ctx)

	defer func() {
		err := app.Scheduler.StopJobs()
		if err != nil {
			log.Printf("Error stopping jobs: %v", err)
		}
		err = app.Scheduler.Shutdown()
		if err != nil {
			log.Fatal(err)
		}
	}()

	q := sqlc.New(app.DB)
	txRunner := database.NewPGXTxRunner(app.DB)
	tracer := otel.Tracer(string(telemetry.UniventsTracer))
	logs := telemetry.Log()

	authMW := middleware.NewAuthMiddleware(app.GaClient, tracer)

	eventRepo := infrastructure.NewEventRepo(q, logs, tracer)
	editionRepo := infrastructure.NewEditionRepo(q, logs, tracer)
	activityRepo := infrastructure.NewActivityRepo(q, logs, tracer)
	checkpointRepo := infrastructure.NewCheckpointRepo(q, logs, tracer)
	ticketRepo := commerceInfra.NewTicketsRepo(q, logs, tracer)
	productRepo := commerceInfra.NewProductsRepo(q, logs, tracer)

	workerHandlers := async.New(editionRepo, app.GaClient, tracer, txRunner)
	activitiesAsyncHandlers := activityAsync.New(activityRepo, app.GaClient, tracer, txRunner)
	server, asynqClient, scheduler, err := worker.InitAsynq(worker.Deps{
		Handlers:         workerHandlers,
		ActivityHandlers: activitiesAsyncHandlers,
	})
	defer func() {
		scheduler.Shutdown()
		server.Shutdown()
		if err = asynqClient.Close(); err != nil {
			telemetry.Log().Error("error closing the asynq client", zap.Error(err))
		}
	}()

	eventCommands := commands.New(eventRepo, app.GaClient, tracer, txRunner)
	eventQueries := queries.New(eventRepo, app.GaClient, tracer, txRunner)
	editionC := editionCommands.New(eventRepo, editionRepo, asynqClient, app.GaClient, tracer, txRunner)
	editionQ := editionQueries.New(eventRepo, editionRepo, app.GaClient, tracer, txRunner)
	activitiesC := activityCommands.New(activityRepo, editionRepo, asynqClient, app.GaClient, tracer, txRunner)
	activitiesQ := activityQueries.New(activityRepo, editionRepo, app.GaClient, tracer, txRunner)
	checkpointsC := checkpointCommands.New(checkpointRepo, editionRepo, asynqClient, app.GaClient, tracer, txRunner)

	ticketsC := ticketsCommands.New(editionRepo, ticketRepo, asynqClient, app.GaClient, tracer, txRunner)
	ticketsQ := ticketsQueries.New(ticketRepo, editionRepo, app.GaClient, tracer, txRunner)
	productsC := productsCommands.New(editionRepo, productRepo, asynqClient, app.GaClient, tracer, txRunner)
	productsQ := productsQueries.New(productRepo, editionRepo, app.GaClient, tracer, txRunner)

	eventHandler := eventhttp.NewEventsHandler(eventCommands, eventQueries)
	editionHandler := editionhttp.NewEditionsHandler(editionC, editionQ)
	activityHandler := activityhttp.NewActivitiesHandler(activitiesC, activitiesQ)
	checkpointHandler := checkpointshttp.NewCheckpointsHandler(checkpointsC)
	ticketHandler := tickethttp.NewTicketsHandler(ticketsC, ticketsQ)
	productHandler := productshttp.NewProductsHandler(productsC, productsQ)

	systemHandler := system.NewUniventsHandler()

	asynqmonHandler := asynqmon.New(asynqmon.Options{
		RootPath: "/admin/asynq",
		RedisConnOpt: asynq.RedisClientOpt{
			Addr:     viper.GetString("REDIS_ADDR"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
	})

	deps := &router.HTTPDeps{
		EventsHandler:      eventHandler,
		EditionsHandler:    editionHandler,
		ActivitiesHandler:  activityHandler,
		CheckpointsHandler: checkpointHandler,
		TicketsHandler:     ticketHandler,
		ProductsHandler:    productHandler,
		SystemHandler:      systemHandler,
		AuthMiddleware:     authMW,
		AsynqmonHandler:    asynqmonHandler,
	}

	app.Deps = deps

	if !skipMux {
		mux := router.CreateRouter(deps)

		log.Printf("GoAuth listening on :%s", app.Port)
		log.Fatal(http.ListenAndServe(":"+app.Port, mux))
	}
}
