package app

import (
	"context"
	"lib/database"
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
	Config    Config
}

var app Informd

func Start() {
	ctx := context.Background()
	SetupConstraintMessages()

	app.Config = LoadConfig()

	SetupFUN(app.Config.AppName)

	app.idxClient = SetupIdentityX(app.Config)

	app.db = database.SetupDB(app.Config.ToDBConfig())
	defer database.CloseDB(app.db)
	app.scheduler = SetupCron(app.db)

	shutdown := telemetry.InitTracer(ctx, app.Config.AppName)
	defer telemetry.ShutdownTracer(ctx, shutdown, app.Config.AppName)

	app.run()
}
