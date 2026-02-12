package initialization

import (
	"context"
	"log"
	"net/http"
	"time"
	"univents/internal/adapters/http/router"
	"univents/internal/infrastructure/telemetry"

	"github.com/TrieOH/goauth-sdk-go"
	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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
	SetupDB(&app, "./internal/database/migrations")
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

	mux, _ := router.CreateRouter(app.gaClient, app.DB, app.Redis)

	log.Printf("GoAuth listening on :%s", app.Port)
	log.Fatal(http.ListenAndServe(":"+app.Port, mux))
}
