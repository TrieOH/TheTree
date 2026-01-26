package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/XSAM/otelsql"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func WaitForDB(timeout time.Duration) (*sql.DB, error) {
	dsn := viper.GetString("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Couldn't get DATABASE_URL variable")
	}

	driverName, err := registerDBDriver()
	if err != nil {
		return nil, err
	}

	deadline := time.Now().Add(timeout)
	attempt := 0

	for {
		attempt++
		log.Printf("waiting for database... attempt %d\n", attempt)

		db, err := sql.Open(driverName, dsn)
		if err != nil {
			log.Printf("error opening the database connection: %v\n", err)
		} else {
			if pingErr := db.Ping(); pingErr == nil {
				log.Printf("database connected on attempt %d\n", attempt)
				return db, nil
			} else {
				log.Printf("error pinging the database: %v\n", pingErr)
				err = db.Close()
				if err != nil {
					log.Printf("error closing the database connection: %v\n", err)
				}
			}
		}

		if time.Now().After(deadline) {
			return nil, errors.New("database connection timeout")
		}

		time.Sleep(2 * time.Second)
	}
}

var (
	otelDriverName string
	registerOnce   sync.Once
)

func registerDBDriver() (string, error) {
	var err error

	registerOnce.Do(func() {
		otelDriverName, err = otelsql.Register(
			"postgres",
			otelsql.WithAttributes(
				semconv.DBSystemPostgreSQL,
			),
			otelsql.WithSQLCommenter(true),
		)
	})

	if err != nil {
		err = fmt.Errorf("error registering otelsql driver: %w", err)
	}

	return otelDriverName, err
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
