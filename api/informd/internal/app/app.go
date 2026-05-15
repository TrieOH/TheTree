package app

import (
	"context"
	"lib/errx"
	"lib/telemetry"

	idx "git.trieoh.com/TrieOH/IdentityX-SDK-Go"
	"github.com/authzed/authzed-go/v1"
	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Informd struct {
	db        *pgxpool.Pool
	redis     *redis.Client
	scheduler gocron.Scheduler
	idxClient *idx.Client
	sdbClient *authzed.Client
	dbErr     *errx.DBHandler
	Config    Config
}

func New() *Informd {
	var app Informd
	var err error
	app.Config, err = LoadConfig()
	if err != nil {
		errx.Must(err, "failed to load config")
	}
	SetupFUN(app.Config.AppName)
	app.idxClient = SetupIdentityX(app.Config)
	migrationPath := "./internal/database/migrations"
	app.dbErr = SetupDBErrorHandler()
	app.db = SetupDB(migrationPath, app.Config, app.dbErr)
	app.scheduler = SetupCron(app.db)
	return &app
}

func (app *Informd) Run() {
	ctx := context.Background()

	defer app.CloseDB()
	defer app.CloseRedis()
	shutdown := app.StartTracer(ctx, app.Config.AppName)
	defer app.ShutdownTracer(ctx, shutdown)
	app.run()
}

func (app *Informd) CloseDB() {
	app.db.Close()
}

func (app *Informd) CloseRedis() {
	if err := app.redis.Close(); err != nil {
		errx.Must(err, "error closing redis connection")
	}
}

func (app *Informd) StartTracer(ctx context.Context, appName string) func(context.Context) error {
	shutdown, err := telemetry.InitTracer(ctx, appName)
	if err != nil {
		errx.Must(err, "error starting tracer")
	}
	return shutdown
}

func (app *Informd) ShutdownTracer(ctx context.Context, shutdown func(context.Context) error) {
	if err := shutdown(ctx); err != nil {
		errx.Must(err, "error shutting down tracer")
	}
}
