// Package testdb spins up disposable Postgres instances backed by
// testcontainers, for integration tests that need a real database.
//
//	pool := testdb.Postgres(t, "../../../../db/migrations")
//
// One container per test binary: the first call in a package starts a
// Postgres container and applies the goose migrations at mPath exactly
// once; every later call in the same binary reuses the same pool. After
// each test the schema's tables are truncated (RESTART IDENTITY CASCADE),
// so tests share the container without seeing each other's rows.
//
// The container lives for the whole test binary and is removed by the
// testcontainers reaper (Ryuk) when the process exits, so there is no
// per-test terminate — starting a container per test was the dominant cost
// of the test suite.
package testdb

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"lib/database"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // register the "pgx" database/sql driver for wait.ForSQL
	"github.com/moby/moby/api/types/network"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Postgres returns a pool connected to the package's shared disposable
// Postgres. The container is started and the migrations at mPath applied
// once per test binary (first call wins); subsequent calls return the same
// pool. A truncate-all runs after each test (and its subtests), so the
// next test starts from a clean schema without paying for a container +
// migrations again.
//
// Pass an empty mPath to skip migrations (e.g. when the test manages its
// own schema).
//
// Tests share the pool, so they must leave no transactions open at the end
// of the test body (pgxpool rolls aborted transactions back on release,
// and the truncate in t.Cleanup would otherwise block on their locks).
//
// When no Docker daemon is reachable the test is SKIPPED rather than
// failed (testcontainers.SkipIfProviderIsNotHealthy): dagger CI runs
// `go test` inside a container with no docker access, so the unit tests
// still run there while the integration tests are executed on the runner
// by the CI workflow (see .forgejo/workflows/ci.yml).
func Postgres(tb testing.TB, mPath string) *pgxpool.Pool {
	tb.Helper()

	if t, ok := tb.(*testing.T); ok {
		testcontainers.SkipIfProviderIsNotHealthy(t)
	}

	ent := entryFor(mPath)
	ent.once.Do(func() {
		ent.pool, ent.err = start(mPath)
	})
	if ent.err != nil {
		tb.Fatalf("testdb: start postgres container: %v", ent.err)
	}

	tb.Cleanup(func() {
		err := truncateAll(ent.pool)
		if err != nil {
			tb.Errorf("testdb: reset schema after test: %v", err)
		}
	})

	return ent.pool
}

type entry struct {
	once sync.Once
	pool *pgxpool.Pool
	err  error
}

var (
	mu      sync.Mutex
	entries = map[string]*entry{}
)

// entryFor returns the per-migration-path singleton. Keying by mPath keeps
// callers with different schema needs (migrated vs raw) in the same binary
// on separate containers instead of silently sharing the wrong schema.
func entryFor(mPath string) *entry {
	mu.Lock()
	defer mu.Unlock()
	e, ok := entries[mPath]
	if !ok {
		e = &entry{}
		entries[mPath] = e
	}
	return e
}

func start(mPath string) (*pgxpool.Pool, error) {
	ctx := context.Background()
	container, err := postgres.Run(ctx, "postgres:18-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testdb"),
		postgres.WithPassword("testdb"),
		testcontainers.WithWaitStrategy(
			// Poll TCP with a real query instead of counting log lines: the
			// docker-entrypoint's temporary initdb server only listens on a
			// unix socket, so the first successful connection is the real
			// server — this returns as soon as the server is usable and
			// avoids the fragile "wait for the second ready log" dance.
			wait.ForSQL("5432/tcp", "pgx", func(host string, port network.Port) string {
				return "postgres://testdb:testdb@" + net.JoinHostPort(host, port.Port()) + "/testdb?sslmode=disable"
			}).WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return nil, err
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	if mPath != "" {
		err := database.RunMigrations(pool, mPath)
		if err != nil {
			return nil, fmt.Errorf("migrations: %w", err)
		}
	}

	return pool, nil
}

// truncateAll empties every table in the public schema, resetting identity
// sequences and cascading through foreign keys, so the next test starts
// from a clean slate while the schema (migrations, types, extensions)
// stays put.
func truncateAll(pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rows, err := pool.Query(ctx, `SELECT quote_ident(tablename) FROM pg_tables WHERE schemaname = 'public'`)
	if err != nil {
		return err
	}
	var tables []string
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			rows.Close()
			return err
		}
		tables = append(tables, name)
	}
	rows.Close()
	err = rows.Err()
	if err != nil {
		return err
	}
	if len(tables) == 0 {
		return nil
	}
	_, err = pool.Exec(ctx, "TRUNCATE "+strings.Join(tables, ", ")+" RESTART IDENTITY CASCADE")
	return err
}
