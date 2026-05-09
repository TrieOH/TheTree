package app

import (
	"IdentityX/internal/platform/telemetry"
	"context"
	"lib/errx"
	telemetry2 "lib/telemetry"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type IdentityX struct {
	db            *pgxpool.Pool
	redis         *redis.Client
	scheduler     gocron.Scheduler
	cfg           Config
	encryptionKey []byte
}

func New() *IdentityX {
	var app IdentityX

	cfg, err := LoadConfig()
	app.cfg = cfg
	if err != nil {
		errx.Must(err, "error loading config")
	}
	app.encryptionKey = InitEncryption(app.cfg.EncryptionKey)
	SetupFUN()
	app.redis = SetupRedis(15*time.Second, app.cfg)
	migrationPath := "./internal/platform/database/migrations"
	app.db = SetupDB(migrationPath, app.cfg.DatabaseURL)
	SetupRuntimeEnv(app.db, app.encryptionKey, app.cfg.KeyLifetime)
	app.scheduler = SetupCron(app.db, app.encryptionKey, app.cfg)

	return &app
}

func (app *IdentityX) Run() {
	ctx := context.Background()

	defer app.CloseDB()
	defer app.CloseRedis()
	defer app.StopScheduler()
	shutdown := app.StartTracer(ctx)
	app.ShutdownTracer(ctx, shutdown)
	app.run()
}

func (app *IdentityX) CloseDB() {
	app.db.Close()
}

func (app *IdentityX) CloseRedis() {
	if err := app.redis.Close(); err != nil {
		telemetry2.Log().Error("error closing redis connection", zap.Error(err))
	}
}

func (app *IdentityX) StartTracer(ctx context.Context) func(context.Context) error {
	shutdown, err := telemetry.InitTracer(ctx)
	if err != nil {
		telemetry2.Log().Fatal("error starting tracer", zap.Error(err))
	}
	return shutdown
}

func (app *IdentityX) ShutdownTracer(ctx context.Context, shutdown func(context.Context) error) {
	if err := shutdown(ctx); err != nil {
		telemetry2.Log().Error("error shutting down tracer", zap.Error(err))
	}
}

func (app *IdentityX) StopScheduler() {
	err := app.scheduler.StopJobs()
	if err != nil {
		telemetry2.Log().Error("error stopping scheduler jobs", zap.Error(err))
	}
	err = app.scheduler.Shutdown()
	if err != nil {
		telemetry2.Log().Error("error shutting down scheduler", zap.Error(err))
	}
}
