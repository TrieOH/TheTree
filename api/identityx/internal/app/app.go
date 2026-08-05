package app

import (
	"IdentityX/internal/config"
	"IdentityX/internal/emails"
	"IdentityX/internal/jobs"
	"IdentityX/internal/repos"
	"IdentityX/internal/sqlc"
	"context"

	"lib/database"
	"lib/email"
	"lib/errx"
	"lib/globals"
	"lib/httpserver"
	libriver "lib/river"
	"lib/telemetry"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
)

type IdentityX struct {
	db          *pgxpool.Pool
	river       *river.Client[pgx.Tx]
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

	libriver.Migrate(ctx, app.db)
	app.river = libriver.NewClient(app.db, libriver.NewWorkers(
		libriver.Register[jobs.CreateCryptoKeyArgs](jobs.NewCreateCryptoKeyWorker(sqlcQueries)),
		libriver.Register[jobs.CleanupBlacklistArgs](jobs.NewCleanupBlacklistWorker(sqlcQueries)),
		libriver.Register[jobs.CleanupActionTokensArgs](jobs.NewCleanupActionTokensWorker(sqlcQueries)),
		libriver.Register[emails.SendAuthEmailArgs](jobs.NewSendAuthEmailWorker(app.emailClient, repos.NewEmailTemplates(sqlcQueries))),
	), nil, []*river.PeriodicJob{
		river.NewPeriodicJob(
			river.PeriodicInterval(5*time.Minute),
			func() (river.JobArgs, *river.InsertOpts) {
				return jobs.CleanupBlacklistArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true},
		),
		river.NewPeriodicJob(
			river.PeriodicInterval(5*time.Minute),
			func() (river.JobArgs, *river.InsertOpts) {
				return jobs.CleanupActionTokensArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true},
		),
	})
	err = app.river.Start(ctx)
	if err != nil {
		errx.Exit(err, "failed to start river client")
	}
	defer func() {
		err = app.river.Stop(ctx)
		if err != nil {
			telemetry.Log().Info("failed to stop river client")
		}
	}()

	EnsureKeysExist(ctx, app.db, app.river)

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
