package app

import (
	"context"
	"lib/database"
	"lib/email"
	"lib/httpserver"
	"lib/objectstorage"
	"lib/telemetry"
	"univents/internal/config"

	idx "sdk/identityx"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Univents struct {
	db          *pgxpool.Pool
	idxClient   *idx.Client
	objStorage  *objectstorage.Client
	emailClient *email.Client

	cfg config.Config
}

var app Univents

func Start() {
	ctx := context.Background()
	SetupConstraintMessages()

	app.cfg = config.Load()

	httpserver.SetupFUN(app.cfg.AppName)

	app.idxClient = SetupIdentityX(app.cfg)
	app.objStorage = SetupObjectStorage(app.cfg)
	app.emailClient = email.NewClient(app.cfg.ToEmailConfig())

	app.db = database.SetupDB(app.cfg.ToDBConfig())
	defer database.CloseDB(app.db)

	shutdown := telemetry.InitTracer(ctx, app.cfg.AppName)
	defer telemetry.ShutdownTracer(ctx, shutdown, app.cfg.AppName)

	app.run()
}
