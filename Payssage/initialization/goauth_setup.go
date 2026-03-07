package initialization

import (
	"TriePayments/internal/interfaces/http/middleware"
	"TriePayments/internal/interfaces/http/router"
	"TriePayments/internal/interfaces/http/system"
	"TriePayments/internal/plataform/telemetry"
	"TriePayments/internal/worker"
	"context"
	"log"
	"net/http"
	"time"

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

type TriePayments struct {
	Port      string
	DB        *pgxpool.Pool
	Redis     *redis.Client
	Scheduler gocron.Scheduler
	GaClient  *goauth.Client

	Deps *router.HTTPDeps
}

func TriePaymentsSetup() *TriePayments {
	var app TriePayments

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

func TriePaymentsStart(app *TriePayments, skipMux bool) {
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

	//q := sqlc.New(app.DB)
	//txRunner := database.NewPGXTxRunner(app.DB)
	tracer := otel.Tracer(string(telemetry.TriePaymentsTracer))
	//logs := telemetry.Log()
	//ws := sockets.New()

	authMW := middleware.NewAuthMiddleware(app.GaClient, tracer)

	// Init Repos

	// init Async Handlers
	server, asynqClient, scheduler, inspector, err := worker.InitAsynq(worker.Deps{})
	defer func() {
		if err = inspector.Close(); err != nil {
			telemetry.Log().Error("error closing the asynq inspector", zap.Error(err))
		}
		scheduler.Shutdown()
		server.Shutdown()
		if err = asynqClient.Close(); err != nil {
			telemetry.Log().Error("error closing the asynq client", zap.Error(err))
		}
	}()

	// Init Commands and Queries

	// Init Handlers

	systemHandler := system.NewSystemHandler()

	asynqmonHandler := asynqmon.New(asynqmon.Options{
		RootPath: "/admin/asynq",
		RedisConnOpt: asynq.RedisClientOpt{
			Addr:     viper.GetString("REDIS_ADDR"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
	})

	deps := &router.HTTPDeps{
		SystemHandler:   systemHandler,
		AuthMiddleware:  authMW,
		AsynqmonHandler: asynqmonHandler,
	}

	app.Deps = deps

	if !skipMux {
		mux := router.CreateRouter(deps)

		log.Printf("TriePayments listening on :%s", app.Port)
		log.Fatal(http.ListenAndServe(":"+app.Port, mux))
	}
}
