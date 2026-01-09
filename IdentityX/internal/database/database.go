package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
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

func SetJWTMasterKey(db *sql.DB) error {
	masterKey := viper.GetString("JWT_MASTER_KEY")
	if masterKey == "" {
		return errors.New("missing JWT_MASTER_KEY in config")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := fmt.Sprintf("SET app.jwt_master_key = '%s'", strings.ReplaceAll(masterKey, "'", "''"))

	_, err := db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to set app.jwt_master_key: %w", err)
	}

	query = fmt.Sprintf("ALTER ROLE CURRENT_USER SET app.jwt_master_key = '%s'", masterKey)
	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to set persistent master key: %w", err)
	}

	log.Println("session variable app.jwt_master_key set successfully")
	return nil
}
