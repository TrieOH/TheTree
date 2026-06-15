package app

import (
	"context"

	"lib/telemetry"

	"github.com/authzed/authzed-go/v1"
	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	idx "sdk/identityx"
)

type Payssage struct {
	db        *pgxpool.Pool
	scheduler gocron.Scheduler
	ga        *idx.Client
	sdb       *authzed.Client
}

func New() *Payssage {
	var app Payssage
	LoadEnv()
	SetupFUN()
	app.ga = SetupIdentityX()
	app.db = SetupDB()
	app.scheduler = SetupCron(app.db)
	app.sdb = SetupSpiceDB()
	return &app
}

func (app *Payssage) Run() {
	ctx := context.Background()

	defer app.CloseDB()
	defer app.StopScheduler()
	shutdown := app.StartTracer(ctx)
	defer app.ShutdownTracer(ctx, shutdown)
	app.run()
}

func (app *Payssage) CloseDB() {
	app.db.Close()
}

func (app *Payssage) StartTracer(ctx context.Context) func(context.Context) error {
	shutdown := telemetry.InitTracer(ctx, "Payssage")
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
