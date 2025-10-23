package main

import (
	"log"
	"time"
	"context"
	"database/sql"

	database "GoAuth/internal/db"
	"GoAuth/internal/repository"

	resp "github.com/MintzyG/GoResponse/response"
	_ "github.com/lib/pq"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/go-co-op/gocron/v2"
)

var Port string
var DB *sql.DB
var scheduler gocron.Scheduler

func init() {
	viper.AutomaticEnv()

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
	})

	var err error
	DB, err = database.WaitForDB(30 * time.Second)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}

	if err := database.RunMigrations(DB, "./migrations"); err != nil {
		log.Fatalf("Failed migrations: %v", err)
	}

	// Create the scheduler
	scheduler, err = gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}

	// Schedule a daily job at 00:00
	_, err = scheduler.NewJob(
		//gocron.DurationJob(1*time.Minute),
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(0, 0, 0))),
		gocron.NewTask(func(ctx context.Context, db *sql.DB) {
			queries := repository.New(db)
			if err := queries.DeleteExpiredTokens(ctx); err != nil {
				log.Printf("Couldn't delete expired tokens: %v\n", err)
			} else {
				log.Println("Expired tokens cleanup executed successfully")
			}
		}, DB),
	)
	if err != nil {
		log.Fatalf("Couldn't create DeleteExpiredTokens cron job: %v", err)
	} else {
		log.Println("Created DeleteExpiredTokens cron job")
	}


	// Schedule a daily job at 00:00
	_, err = scheduler.NewJob(
		//gocron.DurationJob(1*time.Minute),
		gocron.DailyJob(1, gocron.NewAtTimes(gocron.NewAtTime(0, 0, 0))),
		gocron.NewTask(func(ctx context.Context, db *sql.DB) {
			queries := repository.New(db)
			revoked_sessions, err := queries.DeleteExpiredSessions(ctx)
			if err != nil {
				log.Printf("Couldn't delete expired sessions: %v\n", err)
			} else {
				log.Println("Expired sessions cleanup executed successfully")
			}

			tokenIDs := make([]uuid.UUID, len(revoked_sessions))
			expiresAt := make([]time.Time, len(revoked_sessions))

			for i, session := range revoked_sessions {
				tokenIDs[i] = session.TokenID
				expiresAt[i] = session.ExpiresAt
			}

			_, err = queries.BlacklistManyTokens(ctx, repository.BlacklistManyTokensParams{
				Column1: tokenIDs,
				Column2: expiresAt,
			})

			if err != nil {
				log.Println("Couldn't blacklist tokens from invalidated sessions")
			}
		}, DB),
	)
	if err != nil {
		log.Fatalf("Couldn't create DeleteExpiredSessions cron job: %v", err)
	} else {
		log.Println("Created DeleteExpiredSessions cron job")
	}

	// Start the scheduler in the background
	go scheduler.Start()
	log.Println("Started the cron scheduler")
}
