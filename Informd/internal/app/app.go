package app

import (
	"TrieForms/internal/platform/telemetry"
	"context"
	"time"

	"github.com/TrieOH/IdentityX-SDK-Go"
	"github.com/authzed/authzed-go/v1"
	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type TrieForms struct {
	db        *pgxpool.Pool
	redis     *redis.Client
	scheduler gocron.Scheduler
	idxClient *idx.Client
	sdbClient *authzed.Client
}

func New() *TrieForms {
	var app TrieForms

	LoadEnv()
	SetupFUN()
	app.redis = SetupRedis(15 * time.Second)
	app.idxClient = SetupGoAuth()
	migrationPath := "./internal/platform/database/migrations"
	app.db = SetupDB(migrationPath)
	app.scheduler = SetupCron(app.db)
	app.sdbClient = SetupSpiceDB()
	return &app
}

func (app *TrieForms) Run() {
	ctx := context.Background()

	defer app.CloseDB()
	defer app.CloseRedis()
	defer app.StopScheduler()
	shutdown := app.StartTracer(ctx)
	defer app.ShutdownTracer(ctx, shutdown)
	app.run()
}

func (app *TrieForms) CloseDB() {
	app.db.Close()
}

func (app *TrieForms) CloseRedis() {
	if err := app.redis.Close(); err != nil {
		telemetry.Log().Error("error closing redis connection", zap.Error(err))
	}
}

func (app *TrieForms) StartTracer(ctx context.Context) func(context.Context) error {
	shutdown, err := telemetry.InitTracer(ctx)
	if err != nil {
		telemetry.Log().Fatal("error starting tracer", zap.Error(err))
	}
	return shutdown
}

func (app *TrieForms) ShutdownTracer(ctx context.Context, shutdown func(context.Context) error) {
	if err := shutdown(ctx); err != nil {
		telemetry.Log().Error("error shutting down tracer", zap.Error(err))
	}
}

func (app *TrieForms) StopScheduler() {
	err := app.scheduler.StopJobs()
	if err != nil {
		telemetry.Log().Error("error stopping scheduler jobs", zap.Error(err))
	}
	err = app.scheduler.Shutdown()
	if err != nil {
		telemetry.Log().Error("error shutting down scheduler", zap.Error(err))
	}
}
