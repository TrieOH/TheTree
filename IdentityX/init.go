package main

import (
	database "GoAuth/internal/db"
	"GoAuth/internal/repository"
	"context"
	"database/sql"
	resp "github.com/MintzyG/GoResponse/response"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/go-co-op/gocron/v2"
	"log"
	"time"
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
		// gocron.DurationJob(1*time.Minute),
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
		log.Fatalf("Couldn't create cron job: %v", err)
	}

	// Start the scheduler in the background
	go scheduler.Start()
	log.Println("Started the cron scheduler")
}
