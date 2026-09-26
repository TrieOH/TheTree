package config

import (
	"testing"

	"github.com/caarlos0/env/v11"
)

// parse runs env parsing against a fixed environment instead of
// os.Environ, so tests exercise the tags exactly as boot would.
func parse(t *testing.T, environ map[string]string, cfg any) {
	t.Helper()
	err := env.ParseWithOptions(cfg, env.Options{Environment: environ})
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
}

func TestPostgresPrefixInjection(t *testing.T) {
	var cfg struct {
		Postgres Postgres `envPrefix:"UNIVENTS_"`
	}
	parse(t, map[string]string{
		"UNIVENTS_POSTGRES_HOST":     "db-host",
		"UNIVENTS_POSTGRES_DB":       "univents",
		"UNIVENTS_POSTGRES_USER":     "univents",
		"UNIVENTS_POSTGRES_PASSWORD": "secret",
	}, &cfg)

	if cfg.Postgres.Host != "db-host" {
		t.Errorf("Host = %q, want db-host", cfg.Postgres.Host)
	}
	if cfg.Postgres.Port != "5432" {
		t.Errorf("Port default = %q, want 5432", cfg.Postgres.Port)
	}
	if cfg.Postgres.DB != "univents" || cfg.Postgres.User != "univents" || cfg.Postgres.Password != "secret" {
		t.Errorf("Postgres fields not parsed: %+v", cfg.Postgres)
	}
}

func TestRootPostgresDefaults(t *testing.T) {
	var cfg struct {
		RootPostgres
	}
	parse(t, map[string]string{
		"POSTGRES_USER":     "root",
		"POSTGRES_PASSWORD": "root-secret",
	}, &cfg)

	if cfg.User != "root" || cfg.Password != "root-secret" {
		t.Errorf("root creds not parsed: %+v", cfg.RootPostgres)
	}
	if cfg.DB != "postgres" {
		t.Errorf("RootPostgres.DB default = %q, want postgres", cfg.DB)
	}
}

func TestServerDefaults(t *testing.T) {
	var cfg struct {
		Server
	}
	parse(t, map[string]string{"APP_NAME": "svc"}, &cfg)

	if cfg.Port != "8080" {
		t.Errorf("Port default = %q, want 8080", cfg.Port)
	}
	if cfg.AppName != "svc" {
		t.Errorf("AppName = %q, want svc", cfg.AppName)
	}
	if cfg.DebugMode || cfg.DisableRateLimit {
		t.Errorf("flags should default false: %+v", cfg.Server)
	}
}

func TestIdentityXBlock(t *testing.T) {
	var cfg struct {
		IdentityX
	}
	parse(t, map[string]string{
		"IDENTITY_X_URL":        "http://identityx:8080",
		"IDENTITY_X_API_KEY":    "idx_v1_x",
		"IDENTITY_X_PROJECT_ID": "019e4675-3b42-7d39-9a5e-7ba39715a1d3",
	}, &cfg)

	if cfg.URL != "http://identityx:8080" {
		t.Errorf("URL = %q", cfg.URL)
	}
	if cfg.ProjectID.String() != "019e4675-3b42-7d39-9a5e-7ba39715a1d3" {
		t.Errorf("ProjectID = %q", cfg.ProjectID)
	}
}

func TestDBConfigShapesConnections(t *testing.T) {
	got := DBConfig(Postgres{
		Host: "h", Port: "5433", DB: "db", User: "u", Password: "p",
	}, RootPostgres{User: "ru", Password: "rp", DB: "rdb"}, "./db/migrations")

	if got.SSLMode != "disable" || got.RootHost != "postgres" || got.RootPort != "5432" {
		t.Errorf("connection policy drifted: %+v", got)
	}
	if got.Host != "h" || got.Port != "5433" || got.MigrationPath != "./db/migrations" {
		t.Errorf("fields not mapped: %+v", got)
	}
}
