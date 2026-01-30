package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
)

func WaitForDB(timeout time.Duration) (*pgxpool.Pool, error) {
	dsn := viper.GetString("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Couldn't get DATABASE_URL variable")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}

	// Enable OpenTelemetry tracing
	cfg.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Verify connection with retries
	deadline := time.Now().Add(timeout)
	attempt := 0

	for {
		attempt++
		log.Printf("waiting for database... attempt %d\n", attempt)

		if err := pool.Ping(ctx); err == nil {
			log.Printf("database connected on attempt %d\n", attempt)
			return pool, nil
		} else {
			log.Printf("error pinging the database: %v\n", err)
		}

		if time.Now().After(deadline) {
			pool.Close()
			return nil, fmt.Errorf("database connection timeout after %d attempts", attempt)
		}

		time.Sleep(2 * time.Second)
	}
}

// RunMigrations uses pgx/stdlib to provide *sql.DB compatibility for goose
func RunMigrations(pool *pgxpool.Pool, mPath string) error {
	// Convert pgx pool to *sql.DB for goose compatibility
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

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
