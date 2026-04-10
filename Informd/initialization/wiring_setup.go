package initialization

import (
	"TrieForms/internal/features/forms"
	"TrieForms/internal/features/keys"
	"TrieForms/internal/features/projects"
	"TrieForms/internal/interfaces/http/middleware"
	"TrieForms/internal/interfaces/http/router"
	"TrieForms/internal/interfaces/http/system"
	"TrieForms/internal/plataform/database"
	"TrieForms/internal/plataform/database/sqlc"
	"TrieForms/internal/plataform/telemetry"
	"TrieForms/internal/worker"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/TrieOH/goauth-sdk-go"
	v1 "github.com/authzed/authzed-go/v1"
	"github.com/go-co-op/gocron/v2"
	"github.com/hibiken/asynq"
	"github.com/hibiken/asynqmon"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type TrieForms struct {
	Port      string
	DB        *pgxpool.Pool
	Redis     *redis.Client
	Scheduler gocron.Scheduler
	GaClient  *goauth.Client
	AzClient  *v1.Client

	Deps *router.HTTPDeps
}

func TrieFormsSetup() *TrieForms {
	var app TrieForms

	LoadEnv(&app)
	app.Redis = SetupRedis(15 * time.Second)
	SetupSpiceDB(&app)
	SetupGoAuth(&app)
	SetupFUN()
	if viper.GetString("ENV") != "test" {
		SetupDB(&app, "./internal/plataform/database/migrations")
	} else {
		log.Println("WE'RE TESTING")
		SetupDB(&app, "../internal/plataform/database/migrations")
	}
	SetupCron(app.DB, &app)

	return &app
}

func TrieFormsStart(app *TrieForms) {
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
	tracer := otel.Tracer(string(telemetry.TrieFormsTracer))
	logs := telemetry.Log()
	//ws := sockets.New()

	// Init Repos
	projectsRepo := projects.NewProjectRepo(q, logs, tracer)
	apiKeysRepo := keys.NewApiKeyRepo(q, logs, tracer)
	formsRepo := forms.NewFormRepo(q, logs, tracer)

	authMW := middleware.NewAuthMiddleware(app.GaClient, apiKeysRepo, projectsRepo, tracer)

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
	projectsC := projects.NewProjectCommandService(projectsRepo, app.AzClient, txRunner, tracer)
	projectsQ := projects.NewProjectQueryService(projectsRepo, app.AzClient, txRunner, tracer)
	apiKeysC := keys.NewApiKeyCommandService(apiKeysRepo, projectsRepo, app.AzClient, txRunner, tracer)
	apiKeysQ := keys.NewApiKeyQueryService(apiKeysRepo, projectsRepo, app.AzClient, txRunner, tracer)
	formsC := forms.NewFormCommandService(formsRepo, projectsRepo, app.AzClient, txRunner, tracer)
	formsQ := forms.NewFormQueryService(formsRepo, projectsRepo, app.AzClient, txRunner, tracer)

	// Init Handlers
	systemHandler := system.NewSystemHandler(app.GaClient)
	projectsHandler := projects.NewProjectHandler(projectsC, projectsQ)
	apiKeysHandler := keys.NewApiKeysHandler(apiKeysC, apiKeysQ)
	formsHandler := forms.NewFormsHandler(formsC, formsQ)

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
		ProjectsHandler: projectsHandler,
		ApiKeysHandler:  apiKeysHandler,
		FormsHandler:    formsHandler,
		AuthMiddleware:  authMW,
		AsynqmonHandler: asynqmonHandler,
	}

	app.Deps = deps
	mux := router.CreateRouter(deps)

	log.Printf("TrieForms listening on :%s", app.Port)
	log.Fatal(http.ListenAndServe(":"+app.Port, mux))
}
