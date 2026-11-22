package database

import (
	"context"
	"fmt"
	"lib/telemetry"
	"log"
	"log/slog"
	"slices"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func WaitForDB(ctx context.Context, timeout time.Duration, cfg Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	pool, err := tryConnect(ctx, cfg.DSN(), cfg.MaxConns, 5)
	if err != nil {
		return nil, fmt.Errorf("database unreachable: %w", err)
	}
	return pool, nil
}

func CloseDB(pool *pgxpool.Pool) {
	pool.Close()
}

func tryConnect(ctx context.Context, dsn string, maxConns int32, maxAttempts int) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse database URL: %w", err)
	}
	if maxConns > 0 {
		cfg.MaxConns = maxConns
	}
	cfg.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}
	RegisterPoolMetrics(pool)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		telemetry.Log().InfoContext(ctx, "waiting for database", slog.Int("attempt", attempt), slog.Int("max_attempts", maxAttempts))

		err = pingOnce(ctx, pool, maxAttempts-attempt+1)
		if err == nil {
			telemetry.Log().InfoContext(ctx, "database connected", slog.Int("attempt", attempt))
			return pool, nil
		}

		telemetry.Log().WarnContext(ctx, "database ping failed", slog.Int("attempt", attempt), telemetry.Err(err))

		select {
		case <-ctx.Done():
			pool.Close()
			return nil, fmt.Errorf("gave up after %d attempts: %w", attempt, err)
		case <-time.After(2 * time.Second):
		}
	}

	pool.Close()
	return nil, fmt.Errorf("failed to connect after %d attempts: %w", maxAttempts, err)
}

// pingOnce gives one ping an even share of what is left of ctx's budget, so a
// hung first ping can't use it all and leave no room for the retries.
func pingOnce(ctx context.Context, pool *pgxpool.Pool, attemptsLeft int) error {
	pingCtx := ctx
	if deadline, ok := ctx.Deadline(); ok {
		share := time.Until(deadline) / time.Duration(attemptsLeft)
		var cancel context.CancelFunc
		pingCtx, cancel = context.WithTimeout(ctx, share)
		defer cancel()
	}
	return pool.Ping(pingCtx)
}

// RunMigrations uses pgx/stdlib to provide *sql.DB compatibility for goose
func RunMigrations(pool *pgxpool.Pool, mPath string) error {
	db := stdlib.OpenDBFromPool(pool)
	defer func() {
		_ = db.Close()
	}()

	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}
	log.Println("Running migrations...")
	err = goose.Up(db, mPath)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}
	log.Println("Migrations applied successfully")
	return nil
}

// ConstraintErrors maps the schema's constraint and unique-index names to
// the Backend's named error for a violation (ADR-0016).
type ConstraintErrors map[string]error

// constraintRegistry maps constraint/index names to what repos surface on
// violation. SetupDB — the one owner of DB boot — receives it as an
// argument and stores it here for the error handlers, so there is no
// init-order contract with callers: the registry travels as data into the
// one function that owns DB setup, and boot warns when the two drift apart
// (see RegistryDrift).
var constraintRegistry map[string]error

// UseConstraints installs the registry the error handlers read. SetupDB
// calls it at boot; a test that opens its own pool calls it so violations
// surface as the Backend's named errors, as in production.
func UseConstraints(constraints ConstraintErrors) {
	constraintRegistry = constraints
}

// RegistryDrift compares a registry with the schema's constraints and its
// unique indexes that back none, whatever their name. missing are constraints with no named error: a
// violation answers a generic 500 instead of the error it means. unknown
// are entries naming nothing in the schema: a constraint renamed or dropped
// while its entry stayed. River's leftover tables are skipped.
func RegistryDrift(ctx context.Context, db *pgxpool.Pool, constraints ConstraintErrors) (missing, unknown []string, err error) {
	rows, err := db.Query(ctx, `
		SELECT con.conname
		FROM pg_constraint con
		JOIN pg_class rel ON rel.oid = con.conrelid
		JOIN pg_namespace nsp ON nsp.oid = rel.relnamespace
		WHERE nsp.nspname = 'public'
		AND rel.relname != 'goose_db_version'
		AND rel.relname NOT LIKE 'river\_%'
		AND con.contype IN ('u', 'c', 'x')
		UNION
		SELECT idx.relname
		FROM pg_index i
		JOIN pg_class idx ON idx.oid = i.indexrelid
		JOIN pg_class rel ON rel.oid = i.indrelid
		JOIN pg_namespace nsp ON nsp.oid = rel.relnamespace
		WHERE nsp.nspname = 'public'
		AND rel.relname != 'goose_db_version'
		AND rel.relname NOT LIKE 'river\_%'
		AND i.indisunique
		AND NOT i.indisprimary
		AND NOT EXISTS (SELECT 1 FROM pg_constraint con WHERE con.conindid = i.indexrelid)
	`)
	if err != nil {
		return nil, nil, fmt.Errorf("query constraints: %w", err)
	}
	defer rows.Close()

	inSchema := map[string]bool{}
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, nil, fmt.Errorf("scan constraint: %w", err)
		}
		inSchema[name] = true
		if _, ok := constraints[name]; !ok {
			missing = append(missing, name)
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, nil, fmt.Errorf("read constraints: %w", err)
	}
	for name := range constraints {
		if !inSchema[name] {
			unknown = append(unknown, name)
		}
	}
	slices.Sort(missing)
	slices.Sort(unknown)
	return missing, unknown, nil
}
