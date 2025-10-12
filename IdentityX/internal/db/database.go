package db

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
)

func WaitForDB(timeout time.Duration) (*sql.DB, error) {
	dsn := viper.GetString("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Couldn't get DATABASE_URL variable")
	}

	deadline := time.Now().Add(timeout)
	attempt := 1
	for {
		db, err := sql.Open("postgres", dsn)
		if err == nil {
			if pingErr := db.Ping(); pingErr == nil {
				log.Printf("DB connected on attempt %d\n", attempt)
				return db, nil
			} else {
				log.Printf("Ping err: %v\n", pingErr)
			}
			db.Close()
		}

		if time.Now().After(deadline) {
			return nil, errors.New("DB connection timeout")
		}

		log.Printf("Waiting for DB... attempt %d\n", attempt)
		log.Printf("%v\n", err)
		time.Sleep(2 * time.Second)
		attempt++
	}
}

func RunMigrations(db *sql.DB, mPath string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	log.Println("Running migrations...")
	if err := goose.Up(db, mPath); err != nil {
		return err
	}
	log.Println("Migrations applied successfully")
	return nil
}
