package main

import (
	database "GoAuth/internal/db"
	"database/sql"
	resp "github.com/MintzyG/GoResponse/response"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"time"
)

var Port string
var Db *sql.DB

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
	Db, err = database.WaitForDB(30 * time.Second)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}

	if err := database.RunMigrations(Db, "./migrations"); err != nil {
		log.Fatalf("Failed migrations: %v", err)
	}
}
