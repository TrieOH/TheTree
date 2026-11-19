package database

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	"lib/telemetry"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/codes"
)

type Config struct {
	Host           string
	Port           string
	DB             string
	User           string
	Password       string
	SSLMode        string
	ChannelBinding string
	// MaxConns caps the pgx pool; zero keeps pgx's default.
	MaxConns int32
}

func (c Config) DSN() string {
	ssl := c.SSLMode
	if ssl == "" {
		ssl = "require"
	}
	channelBinding := c.ChannelBinding
	if channelBinding == "" {
		channelBinding = "require"
	}
	query := url.Values{}
	query.Set("sslmode", ssl)
	query.Set("channel_binding", channelBinding)
	return fmt.Sprintf(
		"postgres://%s@%s/%s?%s",
		url.UserPassword(c.User, c.Password).String(), net.JoinHostPort(c.Host, c.port()), url.PathEscape(c.DB), query.Encode(),
	)
}

func (c Config) port() string {
	if c.Port == "" {
		return "5432"
	}
	return c.Port
}

// SetupDB is the Postgres adapter's one boot entry point: connect to the
// service database and validate that every constraint has a named error. It
// never migrates (`just migrate <svc>` owns the schema) and returns errors
// instead of exiting, so the caller's boot path owns process failure. Its
// spans join ctx's trace, so a boot traces its connect and constraint check.
func SetupDB(ctx context.Context, cfg Config, constraints ConstraintErrors) (*pgxpool.Pool, error) {
	connectCtx, span := telemetry.StartSpan(ctx, "db.connect")
	db, err := WaitForDB(connectCtx, 30*time.Second, cfg)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "connect")
		span.End()
		return nil, fmt.Errorf("connect: %w", err)
	}
	span.End()
	UseConstraints(constraints)
	checkCtx, span := telemetry.StartSpan(ctx, "db.constraints")
	defer span.End()
	checkCtx, cancel := context.WithTimeout(checkCtx, 10*time.Second)
	defer cancel()
	err = validateConstraints(checkCtx, db)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "constraints")
		CloseDB(db)
		return nil, err
	}
	return db, nil
}
