package main

import (
	http2 "GoAuth/internal/adapters/http/handlers"
	"GoAuth/internal/adapters/observability/logs"
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/apierr"
	"GoAuth/internal/utils"
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"GoAuth/internal/database"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-co-op/gocron/v2"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var Port string
var DB *sql.DB
var scheduler gocron.Scheduler

// init initializes the application.
// It loads environment variables, sets up ED25519 keys, configures the response utility,
// connects to the database, runs migrations, sets the JWT master key, and schedules cron jobs.
func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := http2.LoadProxyConfig(); err != nil {
		log.Fatalf("LoadProxyConfig failed: %v", err.Error())
	}

	if iss := viper.GetString("ISSUER"); iss == "" {
		log.Fatalf("ISSUER environment variable not set.")
	}

	if smtpHost := viper.GetString("SMTP_HOST"); smtpHost == "" {
		log.Fatalf("SMTP_HOST environment variable not set.")
	}
	if smtpPort := viper.GetString("SMTP_PORT"); smtpPort == "" {
		log.Fatalf("SMTP_PORT environment variable not set.")
	}
	if smtpUsername := viper.GetString("SMTP_USERNAME"); smtpUsername == "" {
		log.Fatalf("SMTP_USERNAME environment variable not set.")
	}
	if smtpPassword := viper.GetString("SMTP_PASSWORD"); smtpPassword == "" {
		log.Fatalf("SMTP_PASSWORD environment variable not set.")
	}
	if smtpFrom := viper.GetString("SMTP_FROM"); smtpFrom == "" {
		log.Fatalf("SMTP_FROM environment variable not set.")
	}

	env := viper.GetString("ENV")
	if env != "" && env != "production" {
		apierr.IncludeDebugCauses = true
	}

	err := utils.LoadEd25519Keys(
		viper.GetString("JWT_PRIVATE_KEY"),
		viper.GetString("JWT_PUBLIC_KEY"),
	)

	if err != nil {
		log.Fatal(err)
	}

	Port = viper.GetString("PORT")
	if Port == "" {
		Port = "8080"
	}
	resp.SetConfig(resp.Config{
		MaxTraceSize:         50,
		ResponseSizeLimit:    10 * 1024 * 1024, // 10MB
		MaxInterceptorAmount: 20,
		DefaultContentType:   "application/json",
		EnableSizeValidation: true,
		DefaultModule:        "GoAuth-module",
		ErrorHandler:         apierr.ErrToResp,
	})

	DB, err = database.WaitForDB(30 * time.Second)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}

	if err := database.RunMigrations(DB, "./internal/database/migrations"); err != nil {
		log.Fatalf("Failed migrations: %v", err)
	}

	if err := database.SetJWTMasterKey(DB); err != nil {
		log.Fatal(err)
	}

	// Create the scheduler
	scheduler, err = gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}

	// SessionCleanup scheduled daily at 00:00
	_, err = scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(0, 0, 0))),
		gocron.NewTask(func(ctx context.Context, db *sql.DB) {
			ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			queries := sqlc.New(db)

			revoked, err := queries.RevokeExpiredSessions(ctx)
			if err != nil {
				logs.L().Error("Couldn't revoke expired sessions", zap.Error(err))
				return
			}
			logs.L().Info("Revoked expired sessions", zap.Int("count", len(revoked)))

			deleted, err := queries.DeleteRevokedSessions(ctx)
			if err != nil {
				logs.L().Error("Couldn't delete revoked sessions", zap.Error(err))
				return
			}
			logs.L().Info("Deleted revoked sessions", zap.Int("count", len(deleted)))
		}, DB),
	)

	if err != nil {
		log.Fatalf("Couldn't create SessionCleanup cron job: %v", err)
	} else {
		log.Println("Created SessionCleanup cron job")
	}

	// Start the scheduler in the background
	go scheduler.Start()
	log.Println("Started the cron scheduler")
}
