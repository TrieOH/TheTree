package database

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host           string
	Port           string
	DB             string
	User           string
	Password       string
	SSLMode        string
	ChannelBinding string

	MigrationPath string
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
// service database, run migrations, and validate that every constraint
// carries a registered message. It returns errors instead of exiting so the
// caller's boot path owns process failure.
func SetupDB(cfg Config, messages ConstraintRegistry) (*pgxpool.Pool, error) {
	return setupDB(cfg, messages, true)
}

// SetupDBWithoutMigrations connects to the service database and validates its
// constraints without changing its schema. Use this when schema migrations
// are owned by a separate deployment or developer command.
func SetupDBWithoutMigrations(cfg Config, messages ConstraintRegistry) (*pgxpool.Pool, error) {
	return setupDB(cfg, messages, false)
}

func setupDB(cfg Config, messages ConstraintRegistry, runMigrations bool) (*pgxpool.Pool, error) {
	db, err := WaitForDB(30*time.Second, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	if runMigrations {
		err = RunMigrations(db, cfg.MigrationPath)
		if err != nil {
			CloseDB(db)
			return nil, fmt.Errorf("migrate: %w", err)
		}
	}
	constraintRegistry = messages
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	validateConstraints(ctx, db)
	return db, nil
}
