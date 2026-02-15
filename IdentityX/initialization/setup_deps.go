package initialization

import (
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/database"
	"GoAuth/internal/errx"
	"context"
	"log"
	"time"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail/v3"
	"github.com/MintzyG/fail/v3/plugins/localization"
	"github.com/MintzyG/fail/v3/plugins/tracing/otel"
	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func SetupFail() {
	fail.AllowInternalLogs(true)
	fail.AllowStaticMutations(true, false)
	fail.AllowRuntimePanics(true)

	if err := fail.RegisterTranslator(&errx.HTTPTranslator{}); err != nil {
		log.Fatal(err)
	}

	fail.RegisterMapper(&errx.PGXMapper{})

	fail.SetLocalizer(localization.New())
	fail.SetDefaultLocale("en-US")

	tracerPlugin := otel.New(
		otel.WithTracerName("goauth-service"),
		otel.WithMode(otel.RecordSmart),
		otel.WithStackTrace(),
		otel.WithAttributePrefix("goauth.error"),
	)
	fail.SetTracer(tracerPlugin)
}

func SetupFUN() {
	module := viper.GetString("MODULE")
	if module == "" {
		module = "GoAuth-module"
	}

	resp.SetConfig(resp.Config{
		MaxTraceSize:         50,
		ResponseSizeLimit:    10 * 1024 * 1024,
		MaxInterceptorAmount: 20,
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
		DefaultModule:        module,
		ErrorHandler:         errx.ErrToResp,
	})
}

func SetupDB(app *GoauthApp, migrationPath string) {
	var err error
	db, err := database.WaitForDB(30 * time.Second)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}

	app.DB = db

	if err := database.RunMigrations(db, migrationPath); err != nil {
		log.Fatalf("Failed migrations: %v", err)
	}

	SetupRuntimeEnv(db)
}

func SetupRedis(timeout time.Duration) *redis.Client {
	rdb, err := database.WaitForRedis(timeout)
	if err != nil {
		log.Fatalf("Failed to connect Redis: %v", err)
	}

	return rdb
}

func SetupCron(db *pgxpool.Pool, app *GoauthApp) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}

	app.scheduler = scheduler

	txRunner := transactions.NewTxRunner(db)
	rotateKeysJob(ctx, app, txRunner)
	sessionCleanupJob(ctx, app, txRunner)
	tokenReuseCleanupJob(ctx, app)

	go scheduler.Start()
	log.Println("Started the cron scheduler")
}
