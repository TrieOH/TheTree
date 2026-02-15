package initialization

import (
	"context"
	"log"
	"net/http"
	"time"
	"univents/internal/eventcore/application/commands"
	"univents/internal/eventcore/application/queries"
	"univents/internal/eventcore/infrastructure"
	eventhttp "univents/internal/eventcore/interfaces/http"
	"univents/internal/interfaces/http/middleware"
	"univents/internal/interfaces/http/router"
	"univents/internal/interfaces/http/system"
	"univents/internal/plataform/database"
	"univents/internal/plataform/database/sqlc"
	"univents/internal/plataform/telemetry"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
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
	logs := telemetry.L()

	authMW := middleware.NewAuthMiddleware(app.gaClient, tracer)

	eventRepo := infrastructure.NewEventRepo(q, logs, tracer)

	eventCommands := commands.New(eventRepo, app.gaClient, tracer, txRunner)
	eventQueries := queries.New(eventRepo, app.gaClient, tracer, txRunner)

	eventHandler := eventhttp.NewEventsHandler(eventCommands, eventQueries)

	systemHandler := system.NewUniventsHandler()

	deps := &router.HTTPDeps{
		EventsHandler:  eventHandler,
		SystemHandler:  systemHandler,
		AuthMiddleware: authMW,
	}

	mux := router.CreateRouter(deps)

	log.Printf("GoAuth listening on :%s", app.Port)
	log.Fatal(http.ListenAndServe(":"+app.Port, mux))
}
