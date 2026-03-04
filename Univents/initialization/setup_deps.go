package initialization

import (
	"context"
	"log"
	"time"
	"univents/internal/plataform/database"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/TrieOH/goauth-sdk-go"
	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func SetupGoAuth(app *UniventsApp) {
	projectID := uuid.MustParse(viper.GetString("GO_AUTH_PROJECT_ID"))
	client, err := goauth.NewClient(goauth.Config{
		BaseURL:   viper.GetString("GOAUTH_URL"),
		APIKey:    viper.GetString("GOAUTH_API_KEY"),
		ProjectID: projectID,
	})
	if err != nil {
		log.Fatalf("Error creating goauth client: %s", err.Error())
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = client.Tokens.GetJWKS(ctx, false)
		if err != nil {
			log.Fatalf("error fetching initial JWKS: %s", err.Error())
		}
	}()
	app.GaClient = client
}

func SetupFUN() {
	module := viper.GetString("MODULE")
	if module == "" {
		module = "univents-module"
	}

	resp.SetConfig(resp.Config{
		MaxTraceSize:         50,
		ResponseSizeLimit:    10 * 1024 * 1024,
		MaxInterceptorAmount: 20,
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
		DefaultModule:        module,
	})
}

func SetupDB(app *UniventsApp, migrationPath string) {
	var err error
	db, err := database.WaitForDB(30 * time.Second)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}

	app.DB = db

	if err := database.RunMigrations(db, migrationPath); err != nil {
		log.Fatalf("Failed migrations: %v", err)
	}
}

func SetupRedis(timeout time.Duration) *redis.Client {
	rdb, err := database.WaitForRedis(timeout)
	if err != nil {
		log.Fatalf("Failed to connect Redis: %v", err)
	}

	return rdb
}

func SetupCron(db *pgxpool.Pool, app *UniventsApp) {
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create Scheduler: %v", err)
	}

	app.Scheduler = scheduler

	_ = database.NewPGXTxRunner(db)

	go scheduler.Start()
	log.Println("Started the cron Scheduler")
}
