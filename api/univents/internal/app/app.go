package app

import (
	"context"
	"lib/database"
	"lib/telemetry"

	idx "sdk/identityx"
	"sdk/payssage"

	"github.com/authzed/authzed-go/v1"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
)

type Univents struct {
	db        *pgxpool.Pool
	idxClient *idx.Client
	payssage  *payssage.Client
	minio     *minio.Client
	sdbClient *authzed.Client

	cfg Config
}

var app Univents

func Start() {
	ctx := context.Background()
	SetupConstraintMessages()

	app.cfg = LoadConfig()

	SetupFUN(app.cfg.AppName)

	app.idxClient = SetupIdentityX(app.cfg)
	app.payssage = SetupPayssage(app.cfg)
	app.minio = SetupObjectStorage(app.cfg)

	app.db = database.SetupDB(app.cfg.ToDBConfig())
	defer database.CloseDB(app.db)

	shutdown := telemetry.InitTracer(ctx, app.cfg.AppName)
	defer telemetry.ShutdownTracer(ctx, shutdown, app.cfg.AppName)

	app.run()
}
