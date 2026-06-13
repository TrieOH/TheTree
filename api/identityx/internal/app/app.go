package app

import (
	"context"

	"lib/database"
	"lib/telemetry"

	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IdentityX struct {
	db            *pgxpool.Pool
	scheduler     gocron.Scheduler
	cfg           Config
	encryptionKey []byte
}

var app IdentityX

func Start() {
	ctx := context.Background()
	SetupConstraintMessages()
	app.cfg = LoadConfig()
	app.encryptionKey = InitEncryption()
	SetupFUN()

	app.db = database.SetupDB(app.cfg.ToDBConfig())
	defer database.CloseDB(app.db)

	SetupRuntimeEnv(app.db, app.encryptionKey)
	app.scheduler = SetupCron(app.encryptionKey, app.db, app.cfg)

	shutdown := telemetry.InitTracer(ctx, app.cfg.AppName)
	defer telemetry.ShutdownTracer(ctx, shutdown, app.cfg.AppName)

	app.run()
}
