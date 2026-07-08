package app

import (
	"context"
	"lib/database"
	"lib/objectstorage"
	"lib/telemetry"

	idx "sdk/identityx"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Univents struct {
	db         *pgxpool.Pool
	idxClient  *idx.Client
	objStorage *objectstorage.Client

	cfg Config
}

var app Univents

func Start() {
	ctx := context.Background()
	SetupConstraintMessages()

	app.cfg = LoadConfig()

	SetupFUN(app.cfg.AppName)

	app.idxClient = SetupIdentityX(app.cfg)
	app.objStorage = SetupObjectStorage(app.cfg)

	app.db = database.SetupDB(app.cfg.ToDBConfig())
	defer database.CloseDB(app.db)

	shutdown := telemetry.InitTracer(ctx, app.cfg.AppName)
	defer telemetry.ShutdownTracer(ctx, shutdown, app.cfg.AppName)

	app.run()
}
