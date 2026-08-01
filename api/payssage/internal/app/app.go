package app

import (
	"context"
	"lib/database"
	"lib/httpserver"
	libriver "lib/river"
	"lib/telemetry"
	"payssage/internal/config"

	idx "sdk/identityx"

	"github.com/jackc/pgx/v5/pgxpool"
	"resty.dev/v3"
)

type Payssage struct {
	db         *pgxpool.Pool
	idxClient  *idx.Client
	httpClient *resty.Client

	cfg config.Config
}

var app Payssage

func Start() {
	ctx := context.Background()
	SetupConstraintMessages()

	app.cfg = config.LoadConfig()

	httpserver.SetupFUN(app.cfg.AppName)

	app.idxClient = SetupIdentityX(app.cfg)

	app.httpClient = SetupHTTPClient()

	app.db = database.SetupDB(app.cfg.ToDBConfig())
	defer database.CloseDB(app.db)

	libriver.Migrate(ctx, app.db)

	shutdown := telemetry.InitTracer(ctx, app.cfg.AppName)
	defer telemetry.ShutdownTracer(ctx, shutdown, app.cfg.AppName)

	app.run()
}
