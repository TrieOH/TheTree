package app

import (
	"lib/database"
	"lib/errx"

	"github.com/caarlos0/env/v11"
	"github.com/google/uuid"
)

type Config struct {
	// Server
	Port        string `env:"PORT" envDefault:"8080"`
	ProfilePort string `env:"PROFILE_PORT" envDefault:"6060"`
	AppName     string `env:"APP_NAME,required"`
	//AppUrl      string `env:"APP_URL,required"`
	DebugMode  bool   `env:"DEBUG_MODE"`
	HmacSecret string `env:"HMAC_SECRET,required"`

	// Security
	//WsJwtSecret string `env:"WS_JWT_SECRET,required"`

	// IdentityX
	IdxURL       string    `env:"IDENTITY_X_URL,required"`
	IdxAPIKey    string    `env:"IDENTITY_X_API_KEY,required"`
	IdxProjectID uuid.UUID `env:"IDENTITY_X_PROJECT_ID,required"`

	// Payssage
	//PayssageURL           string `env:"PAYSSAGE_URL,required"`
	//PayssageProvider      string `env:"PAYSSAGE_PROVIDER,required"`
	//PayssageAPIKey        string `env:"PAYSSAGE_API_KEY,required"`
	//PayssageWebhookSecret string `env:"PAYSSAGE_WEBHOOK_SECRET,required"`

	// Postgres (own DB)
	PostgresHost     string `env:"UNIVENTS_POSTGRES_HOST,required"`
	PostgresPort     string `env:"UNIVENTS_POSTGRES_PORT" envDefault:"5432"`
	PostgresDB       string `env:"UNIVENTS_POSTGRES_DB,required"`
	PostgresUser     string `env:"UNIVENTS_POSTGRES_USER,required"`
	PostgresPassword string `env:"UNIVENTS_POSTGRES_PASSWORD,required"`

	// Migration
	MigrationPath string `env:"MIGRATION_PATH,required" envDefault:"./internal/database/migrations"`

	// Postgres (root — from .env)
	RootPostgresUser     string `env:"POSTGRES_USER,required"`
	RootPostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	RootPostgresDB       string `env:"POSTGRES_DB" envDefault:"postgres"`

	// Object Storage (RustFS)
	ObjStorageEndpoint  string `env:"OBJECT_STORAGE_ENDPOINT,required"`
	ObjStorageAccessKey string `env:"OBJECT_STORAGE_ACCESS_KEY,required"`
	ObjStorageSecretKey string `env:"OBJECT_STORAGE_SECRET_KEY,required"`
	ObjStorageUseSSL    bool   `env:"OBJECT_STORAGE_USE_SSL" envDefault:"true"`
	ObjStorageRegion    string `env:"OBJECT_STORAGE_REGION"  envDefault:"us-east-1"`

	// CORS
	CorsAllowedOrigins string `env:"CORS_ALLOWED_ORIGINS,required"`
	CorsAllowedHeaders string `env:"CORS_ALLOWED_HEADERS,required"`

	// Feature flags
	DisableRateLimit bool `env:"DISABLE_RATE_LIMIT"`
}

func (cfg Config) ToDBConfig() database.Config {
	return database.Config{
		Host:          cfg.PostgresHost,
		Port:          cfg.PostgresPort,
		DB:            cfg.PostgresDB,
		User:          cfg.PostgresUser,
		Password:      cfg.PostgresPassword,
		SSLMode:       "disable",
		RootUser:      cfg.RootPostgresUser,
		RootPassword:  cfg.RootPostgresPassword,
		RootDB:        cfg.RootPostgresDB,
		RootHost:      "postgres",
		RootPort:      "5432",
		MigrationPath: cfg.MigrationPath,
	}
}

func LoadConfig() Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		errx.Exit(err, "error loading config")
	}
	return cfg
}
