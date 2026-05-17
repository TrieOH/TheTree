package database

import (
	"context"
	"errors"
	"fmt"
	"lib/errx"
	"log"
	"strings"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func WaitForDB(timeout time.Duration, cfg Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := provisionDB(ctx, cfg); err != nil {
		return nil, fmt.Errorf("failed to provision database: %w", err)
	}

	pool, err := tryConnect(ctx, cfg.DSN(), 5)
	if err != nil {
		return nil, fmt.Errorf("database unreachable after provisioning: %w", err)
	}
	return pool, nil
}

func CloseDB(pool *pgxpool.Pool) {
	pool.Close()
}

func tryConnect(ctx context.Context, dsn string, maxAttempts int) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}
	cfg.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		log.Printf("waiting for database... attempt %d/%d\n", attempt, maxAttempts)

		if err = pool.Ping(ctx); err == nil {
			log.Printf("database connected on attempt %d\n", attempt)
			return pool, nil
		}

		log.Printf("ping failed: %v\n", err)

		select {
		case <-ctx.Done():
			pool.Close()
			return nil, fmt.Errorf("context deadline exceeded after %d attempts", attempt)
		case <-time.After(2 * time.Second):
		}
	}

	pool.Close()
	return nil, fmt.Errorf("failed to connect after %d attempts", maxAttempts)
}

func provisionDB(ctx context.Context, cfg Config) error {
	conn, err := pgx.Connect(ctx, cfg.RootDSN())
	if err != nil {
		return fmt.Errorf("unable to connect to root postgres: %w", err)
	}
	defer conn.Close(ctx)

	_, err = conn.Exec(ctx, fmt.Sprintf(
		`DO $$ BEGIN
            IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = '%s') THEN
                CREATE ROLE "%s" LOGIN PASSWORD '%s';
            END IF;
        END $$;`,
		cfg.User, cfg.User, cfg.Password,
	))
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}
	log.Printf("role %q ensured\n", cfg.User)

	var exists bool
	if err = conn.QueryRow(ctx,
		`SELECT EXISTS(SELECT FROM pg_database WHERE datname = $1)`, cfg.DB,
	).Scan(&exists); err != nil {
		return fmt.Errorf("failed to check database existence: %w", err)
	}

	if !exists {
		if _, err = conn.Exec(ctx, fmt.Sprintf(
			`CREATE DATABASE "%s" OWNER "%s"`, cfg.DB, cfg.User,
		)); err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		log.Printf("database %q created\n", cfg.DB)
	} else {
		log.Printf("database %q already exists\n", cfg.DB)
	}

	return nil
}

// RunMigrations uses pgx/stdlib to provide *sql.DB compatibility for goose
func RunMigrations(pool *pgxpool.Pool, mPath string) error {
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}
	log.Println("Running migrations...")
	if err := goose.Up(db, mPath); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	log.Println("Migrations applied successfully")
	return nil
}

type ConstraintRegistry map[string]string

var ConstraintErrorRegistry ConstraintRegistry

func SetConstraintErrorRegistry(registry ConstraintRegistry) {
	ConstraintErrorRegistry = registry
}

func validateConstraints(ctx context.Context, db *pgxpool.Pool) {
	rows, err := db.Query(ctx, `
		SELECT con.conname
		FROM pg_constraint con
		JOIN pg_class rel ON rel.oid = con.conrelid
		JOIN pg_namespace nsp ON nsp.oid = rel.relnamespace
		WHERE nsp.nspname = 'public'
		AND rel.relname != 'goose_db_version'
		AND con.contype IN ('u', 'c')
		UNION
		SELECT indexname
		FROM pg_indexes
		WHERE schemaname = 'public'
		AND tablename != 'goose_db_version'
		AND (indexname LIKE 'uniq_%' OR indexname LIKE 'one_%')
	`)
	if err != nil {
		errx.Exit(err, "error querying constraints")
	}
	defer rows.Close()

	var missing []string
	for rows.Next() {
		var name string
		_ = rows.Scan(&name)
		if _, ok := ConstraintErrorRegistry[name]; !ok {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		errx.Exit(errors.New("missing constraints messages"), "the following constraint messages are missing from the registry: "+strings.Join(missing, ", "))
	}
}
