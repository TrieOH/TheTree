package database

import "fmt"

type Config struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string

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
		c.User, c.Password, c.Host, c.port(), c.Name, ssl,
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
