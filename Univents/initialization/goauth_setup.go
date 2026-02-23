package initialization

import (
	"context"
	"log"
	"net/http"
	"time"
	"univents/internal/core/application/edition/async"
	editionCommands "univents/internal/core/application/edition/commands"
	editionQueries "univents/internal/core/application/edition/queries"
	"univents/internal/core/application/event/commands"
	"univents/internal/core/application/event/queries"
	"univents/internal/core/infrastructure"
	eventhttp "univents/internal/core/interfaces/http"
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
	scheduler gocron.Scheduler
	gaClient  *goauth.Client
}

func UniventsSetup() *UniventsApp {
	var app UniventsApp

	LoadEnv(&app)
	SetupGoAuth(&app)
	SetupFail()
	SetupFUN()
	SetupDB(&app, "./internal/plataform/database/migrations")
	app.Redis = SetupRedis(15 * time.Second)
	SetupCron(app.DB, &app)

	return &app
}

func UniventsStart(app *UniventsApp) {
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
		err := app.scheduler.StopJobs()
		if err != nil {
			log.Printf("Error stopping jobs: %v", err)
		}
		err = app.scheduler.Shutdown()
		if err != nil {
			log.Fatal(err)
		}
	}()

	q := sqlc.New(app.DB)
	txRunner := database.NewPGXTxRunner(app.DB)
	tracer := otel.Tracer(string(telemetry.UniventsTracer))
	logs := telemetry.Log()

	authMW := middleware.NewAuthMiddleware(app.gaClient, tracer)

	eventRepo := infrastructure.NewEventRepo(q, logs, tracer)
	editionRepo := infrastructure.NewEditionRepo(q, logs, tracer)

	workerHandlers := async.New(editionRepo, app.gaClient, tracer, txRunner)
	server, asynqClient, scheduler, err := worker.InitAsynq(worker.Deps{
		Handlers: workerHandlers,
	})
	defer func() {
		scheduler.Shutdown()
		server.Shutdown()
		if err = asynqClient.Close(); err != nil {
			telemetry.Log().Error("error closing the asynq client", zap.Error(err))
		}
	}()

	eventCommands := commands.New(eventRepo, app.gaClient, tracer, txRunner)
	eventQueries := queries.New(eventRepo, app.gaClient, tracer, txRunner)
	editionC := editionCommands.New(eventRepo, editionRepo, asynqClient, app.gaClient, tracer, txRunner)
	editionQ := editionQueries.New(eventRepo, editionRepo, app.gaClient, tracer, txRunner)

	eventHandler := eventhttp.NewEventsHandler(eventCommands, eventQueries)
	editionHandler := editionhttp.NewEditionsHandler(editionC, editionQ)

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
		EventsHandler:   eventHandler,
		EditionsHandler: editionHandler,
		SystemHandler:   systemHandler,
		AuthMiddleware:  authMW,
		AsynqmonHandler: asynqmonHandler,
	}

	mux := router.CreateRouter(deps)

	log.Printf("GoAuth listening on :%s", app.Port)
	log.Fatal(http.ListenAndServe(":"+app.Port, mux))
}
