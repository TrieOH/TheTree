package app

import (
	"context"
	"payssage/internal/platform/telemetry"
	"time"

	"github.com/TrieOH/IdentityX-SDK-Go"
	"github.com/authzed/authzed-go/v1"
	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Payssage struct {
	db        *pgxpool.Pool
	redis     *redis.Client
	scheduler gocron.Scheduler
	ga        *idx.Client
	sdb       *authzed.Client
}

func New() *Payssage {
	var app Payssage
	LoadEnv()
	SetupFUN()
	app.redis = SetupRedis(15 * time.Second)
	app.ga = SetupIdentityX()
	migrationPath := "./internal/platform/database/migrations"
	app.db = SetupDB(migrationPath)
	app.scheduler = SetupCron(app.db)
	app.sdb = SetupSpiceDB()
	return &app
}

func (app *Payssage) Run() {
	ctx := context.Background()

	defer app.CloseDB()
	defer app.CloseRedis()
	defer app.StopScheduler()
	shutdown := app.StartTracer(ctx)
	defer app.ShutdownTracer(ctx, shutdown)
	app.run()
}

func (app *Payssage) CloseDB() {
	app.db.Close()
}

func (app *Payssage) CloseRedis() {
	if err := app.redis.Close(); err != nil {
		telemetry.Log().Error("error closing redis connection", zap.Error(err))
	}
}

func (app *Payssage) StartTracer(ctx context.Context) func(context.Context) error {
	shutdown, err := telemetry.InitTracer(ctx)
	if err != nil {
		telemetry.Log().Fatal("error starting tracer", zap.Error(err))
	}
	return shutdown
}

func (app *Payssage) ShutdownTracer(ctx context.Context, shutdown func(context.Context) error) {
	if err := shutdown(ctx); err != nil {
		telemetry.Log().Error("error shutting down tracer", zap.Error(err))
	}
}

func (app *Payssage) StopScheduler() {
	err := app.scheduler.StopJobs()
	if err != nil {
		telemetry.Log().Error("error stopping scheduler jobs", zap.Error(err))
	}
	err = app.scheduler.Shutdown()
	if err != nil {
		telemetry.Log().Error("error shutting down scheduler", zap.Error(err))
	}
}
