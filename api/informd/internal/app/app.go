package app

import (
	"Informd/internal/config"
	"context"

	"lib/database"
	"lib/httpserver"
	"lib/telemetry"

	idx "sdk/identityx"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Informd struct {
	db        *pgxpool.Pool
	idxClient *idx.Client
	cfg       config.Config
}

var app Informd

func Start() {
	ctx := context.Background()
	SetupConstraintMessages()

	app.cfg = config.LoadConfig()

	httpserver.SetupFUN(app.cfg.AppName)

	app.idxClient = SetupIdentityX(app.cfg)

	app.db = database.SetupDB(app.cfg.ToDBConfig())
	defer database.CloseDB(app.db)

	shutdown := telemetry.InitTracer(ctx, app.cfg.AppName)
	defer telemetry.ShutdownTracer(ctx, shutdown, app.cfg.AppName)

	app.run()
}
