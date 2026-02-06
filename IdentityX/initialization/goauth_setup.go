package initialization

import (
	"GoAuth/internal/adapters/http/router"
	"GoAuth/internal/infrastructure/telemetry"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type GoauthApp struct {
	Port      string
	DB        *pgxpool.Pool
	Redis     *redis.Client
	scheduler gocron.Scheduler
}

func GoAuthSetup() *GoauthApp {
	var app GoauthApp

	LoadEnv(&app)
	SetupFail()
	SetupFUN()
	SetupDB(&app, "./internal/database/migrations")
	app.Redis = SetupRedis(15 * time.Second)
	SetupCron(app.DB, &app)

	return &app
}

func GoAuthStart(app *GoauthApp) {
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

	mux, _ := router.CreateRouter(app.DB, app.Redis)

	log.Printf("GoAuth listening on :%s", app.Port)
	log.Fatal(http.ListenAndServe(":"+app.Port, mux))
}
