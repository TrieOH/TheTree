// Package testdb spins up disposable Postgres instances backed by
// testcontainers, for integration tests that need a real database.
//
//	pool := testdb.Postgres(t, "../../../../db/migrations")
//	defer func() { _ = pool.Close() }() // also auto-terminated via t.Cleanup
package testdb

import (
	"context"
	"testing"
	"time"

	"lib/database"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Postgres starts a disposable Postgres container, applies the goose
// migrations at mPath, and returns a connected pool ready for queries.
// The container is terminated and the pool closed when the test finishes
// (t.Cleanup), so callers only need to use the pool.
//
// Pass an empty mPath to skip migrations (e.g. when the test manages its
// own schema).
func Postgres(t testing.TB, mPath string) *pgxpool.Pool {
	t.Helper()

	ctx := context.Background()
	container, err := postgres.Run(ctx, "postgres:18-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testdb"),
		postgres.WithPassword("testdb"),
		testcontainers.WithWaitStrategy(
			// the module has no default wait strategy: listening port is not
			// enough, and postgres logs "ready" once from its temporary
			// initdb server before restarting and killing connections — wait
			// for the second occurrence, from the real server
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("testdb: start postgres container: %v", err)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("testdb: connection string: %v", err)
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("testdb: open pool: %v", err)
	}
	if err := pool.Ping(ctx); err != nil {
		t.Fatalf("testdb: ping: %v", err)
	}

	if mPath != "" {
		err := database.RunMigrations(pool, mPath)
		if err != nil {
			t.Fatalf("testdb: migrations: %v", err)
		}
	}

	t.Cleanup(func() {
		pool.Close()
		_ = testcontainers.TerminateContainer(container)
	})

	return pool
}
