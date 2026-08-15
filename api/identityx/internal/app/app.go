package app

import (
	"IdentityX/internal/config"
	"IdentityX/internal/sqlc"
	"context"

	"lib/database"
	"lib/email"
	"lib/errx"
	"lib/globals"
	"lib/httpserver"
	"lib/telemetry"

	"github.com/jackc/pgx/v5/pgxpool"
)

type IdentityX struct {
	db          *pgxpool.Pool
	cfg         config.Config
	emailClient *email.Client
}

var app IdentityX

func Start() {
	ctx := context.Background()
	SetupConstraintMessages()
	app.cfg = config.LoadConfig()
	httpserver.SetupFUN(app.cfg.AppName)

	app.db = database.SetupDB(app.cfg.ToDBConfig())
	defer database.CloseDB(app.db)

	app.emailClient = email.NewClient(app.cfg.ToEmailConfig())

	sqlcQueries := sqlc.New(app.db)
	has, err := sqlcQueries.HasAnyActor(ctx)
	if err != nil {
		errx.Exit(err, "failed to check setup state")
	}
	if has {
		globals.MarkSetupComplete()
	}

	shutdown := telemetry.InitTracer(ctx, app.cfg.AppName)
	defer telemetry.ShutdownTracer(ctx, shutdown, app.cfg.AppName)

	app.run()
}
