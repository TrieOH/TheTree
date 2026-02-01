package initialization

import (
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/apierr"
	"GoAuth/internal/database"
	"context"
	"log"
	"time"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail"
	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func SetupFail() {
	fail.AllowInternalLogs(false)

	if err := fail.RegisterTranslator(&apierr.HTTPTranslator{}); err != nil {
		log.Fatal(err)
	}

	fail.RegisterMapper(&apierr.PGXMapper{})

	fail.OnFromFail(apierr.OnFromFailHook)
	fail.OnFromSuccess(apierr.OnFromSuccessHook)
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
		ErrorHandler:         apierr.ErrToResp,
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

	go scheduler.Start()
	log.Println("Started the cron scheduler")
}
