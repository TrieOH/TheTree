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
	Host     string
	Port     string
	DB       string
	User     string
	Password string
	SSLMode  string

	MigrationPath string

	RootUser     string
	RootPassword string
	RootDB       string
	RootHost     string
	RootPort     string
}

func (c Config) DSN() string {
	ssl := c.SSLMode
	if ssl == "" {
		ssl = "disable"
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		url.QueryEscape(c.User), url.QueryEscape(c.Password), net.JoinHostPort(c.Host, c.port()), c.DB, ssl,
	)
}

func (c Config) RootDSN() string {
	host := c.RootHost
	if host == "" {
		host = c.Host
	}
	rootPort := c.RootPort
	if rootPort == "" {
		rootPort = c.port()
	}
	rootDB := c.RootDB
	if rootDB == "" {
		rootDB = "postgres"
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		url.QueryEscape(c.RootUser), url.QueryEscape(c.RootPassword), net.JoinHostPort(host, rootPort), rootDB,
	)
}

func (c Config) port() string {
	if c.Port == "" {
		return "5432"
	}
	return c.Port
}

// SetupDB is the Postgres adapter's one boot entry point: provision the
// role and database, wait for connectivity, run migrations, and validate
// that every constraint carries a registered message. It returns errors
// instead of exiting so the caller's boot path owns process failure.
func SetupDB(cfg Config, messages ConstraintRegistry) (*pgxpool.Pool, error) {
	db, err := WaitForDB(30*time.Second, cfg)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	err = RunMigrations(db, cfg.MigrationPath)
	if err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	constraintRegistry = messages
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	validateConstraints(ctx, db)
	return db, nil
}
