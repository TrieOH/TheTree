package app

import (
	"IdentityX/internal/platform/telemetry"
	"IdentityX/internal/shared/errx"
	"context"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type IdentityX struct {
	db        *pgxpool.Pool
	redis     *redis.Client
	scheduler gocron.Scheduler
	cfg       Config
}

func New() *IdentityX {
	var app IdentityX

	cfg, err := LoadConfig()
	if err != nil {
		errx.Must(err, "error loading config")
	}
	app.cfg = cfg
	SetupFUN()
	app.redis = SetupRedis(15 * time.Second)
	migrationPath := "./internal/platform/database/migrations"
	app.db = SetupDB(migrationPath)
	SetupRuntimeEnv(app.db)
	app.scheduler = SetupCron(app.db)

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
		telemetry.Log().Error("error closing redis connection", zap.Error(err))
	}
}

func (app *IdentityX) StartTracer(ctx context.Context) func(context.Context) error {
	shutdown, err := telemetry.InitTracer(ctx)
	if err != nil {
		telemetry.Log().Fatal("error starting tracer", zap.Error(err))
	}
	return shutdown
}

func (app *IdentityX) ShutdownTracer(ctx context.Context, shutdown func(context.Context) error) {
	if err := shutdown(ctx); err != nil {
		telemetry.Log().Error("error shutting down tracer", zap.Error(err))
	}
}

func (app *IdentityX) StopScheduler() {
	err := app.scheduler.StopJobs()
	if err != nil {
		telemetry.Log().Error("error stopping scheduler jobs", zap.Error(err))
	}
	err = app.scheduler.Shutdown()
	if err != nil {
		telemetry.Log().Error("error shutting down scheduler", zap.Error(err))
	}
}
