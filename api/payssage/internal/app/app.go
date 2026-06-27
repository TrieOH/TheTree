package app

import (
	"context"
	"lib/database"
	libriver "lib/river"
	"lib/telemetry"

	idx "sdk/identityx"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
)

type Payssage struct {
	db        *pgxpool.Pool
	idxClient *idx.Client
	river     *river.Client[pgx.Tx]

	cfg Config
}

var app Payssage

func Start() {
	ctx := context.Background()
	SetupConstraintMessages()

	app.cfg = LoadConfig()

	SetupFUN(app.cfg.AppName)

	app.idxClient = SetupIdentityX(app.cfg)

	app.db = database.SetupDB(app.cfg.ToDBConfig())
	defer database.CloseDB(app.db)

	libriver.Migrate(ctx, app.db)

	shutdown := telemetry.InitTracer(ctx, app.cfg.AppName)
	defer telemetry.ShutdownTracer(ctx, shutdown, app.cfg.AppName)

	app.run()
}
