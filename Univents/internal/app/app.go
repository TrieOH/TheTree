package app

import (
	"context"
	"time"
	"univents/internal/platform/telemetry"

	paymentsSDK "github.com/TrieOH/TriePaymentsSDK"
	"github.com/TrieOH/goauth-sdk-go"
	"github.com/authzed/authzed-go/v1"
	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Univents struct {
	db        *pgxpool.Pool
	redis     *redis.Client
	scheduler gocron.Scheduler
	gaClient  *goauth.Client
	payssage  *paymentsSDK.Client
	minio     *minio.Client
	sdbClient *authzed.Client
}

func New() *Univents {
	var app Univents
	LoadEnv()
	SetupFUN()
	app.redis = SetupRedis(15 * time.Second)
	app.gaClient = SetupGoAuth(app.redis)
	app.payssage = SetupPayssage()
	app.minio = SetupObjectStorage()
	migrationPath := "./internal/platform/database/migrations"
	app.db = SetupDB(migrationPath)
	app.scheduler = SetupCron(&app)
	app.sdbClient = SetupSpiceDB()

	return &app
}

func (app *Univents) Run() {
	ctx := context.Background()

	defer app.CloseDB()
	defer app.CloseRedis()
	defer app.StopScheduler()
	shutdown := app.StartTracer(ctx)
	defer app.ShutdownTracer(ctx, shutdown)
	app.run()
}

func (app *Univents) CloseDB() {
	app.db.Close()
}

func (app *Univents) CloseRedis() {
	if err := app.redis.Close(); err != nil {
		telemetry.Log().Error("error closing redis connection", zap.Error(err))
	}
}

func (app *Univents) StartTracer(ctx context.Context) func(context.Context) error {
	shutdown, err := telemetry.InitTracer(ctx)
	if err != nil {
		telemetry.Log().Fatal("error starting tracer", zap.Error(err))
	}
	return shutdown
}

func (app *Univents) ShutdownTracer(ctx context.Context, shutdown func(context.Context) error) {
	if err := shutdown(ctx); err != nil {
		telemetry.Log().Error("error shutting down tracer", zap.Error(err))
	}
}

func (app *Univents) StopScheduler() {
	err := app.scheduler.StopJobs()
	if err != nil {
		telemetry.Log().Error("error stopping scheduler jobs", zap.Error(err))
	}
	err = app.scheduler.Shutdown()
	if err != nil {
		telemetry.Log().Error("error shutting down scheduler", zap.Error(err))
	}
}
