package database

import (
	"context"
	"fmt"
	"lib/errx"
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
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.port(), c.DB, ssl,
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
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.RootUser, c.RootPassword, host, rootPort, rootDB,
	)
}

func (c Config) port() string {
	if c.Port == "" {
		return "5432"
	}
	return c.Port
}

func SetupDB(cfg Config) *pgxpool.Pool {
	db, err := WaitForDB(30*time.Second, cfg)
	if err != nil {
		errx.Exit(err, "Failed to connect DB")
	}
	if err = RunMigrations(db, cfg.MigrationPath); err != nil {
		errx.Exit(err, "Failed migrations")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	validateConstraints(ctx, db)
	return db
}
