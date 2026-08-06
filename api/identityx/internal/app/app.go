package app

import (
	"IdentityX/internal/config"
	"IdentityX/internal/jobs"
	"IdentityX/internal/sqlc"
	"context"

	"lib/database"
	"lib/email"
	"lib/errx"
	"lib/globals"
	"lib/httpserver"
	"lib/telemetry"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
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

func EnsureKeysExist(ctx context.Context, db *pgxpool.Pool, riverClient *river.Client[pgx.Tx]) {
	q := sqlc.New(db)
	for _, keyType := range []string{"signing", "encryption"} {
		exists, err := q.HasActiveCryptoKey(ctx, sqlc.HasActiveCryptoKeyParams{Type: keyType})
		if err != nil {
			errx.Exit(err, "failed to check global "+keyType+" key")
		}
		if !exists {
			_, err = riverClient.Insert(ctx, jobs.CreateCryptoKeyArgs{KeyType: keyType}, nil)
			if err != nil {
				errx.Exit(err, "failed to enqueue global "+keyType+" key creation")
			}
		}
	}

	projects, err := q.ListProjects(ctx)
	if err != nil {
		errx.Exit(err, "failed to list projects for key check")
	}

	for _, p := range projects {
		pid := p.ID
		for _, keyType := range []string{"signing", "encryption"} {
			exists, err := q.HasActiveCryptoKey(ctx, sqlc.HasActiveCryptoKeyParams{
				ProjectID: &pid,
				Type:      keyType,
			})
			if err != nil {
				errx.Exit(err, "failed to check "+keyType+" key for project "+pid.String())
			}
			if !exists {
				_, err = riverClient.Insert(ctx, jobs.CreateCryptoKeyArgs{
					ProjectID: &pid,
					KeyType:   keyType,
				}, nil)
				if err != nil {
					errx.Exit(err, "failed to enqueue "+keyType+" key creation for project "+pid.String())
				}
			}
		}
	}
}
