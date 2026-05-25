package app

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/internal/jobs"
	"context"
	"lib/database"
	"lib/errx"
	"lib/globals"
	libriver "lib/river"
	"lib/telemetry"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
)

type IdentityX struct {
	db        *pgxpool.Pool
	scheduler gocron.Scheduler
	river     *river.Client[pgx.Tx]
	cfg       Config
}

var app IdentityX

func Start() {
	ctx := context.Background()
	SetupConstraintMessages()
	app.cfg = LoadConfig()
	SetupFUN()

	app.db = database.SetupDB(app.cfg.ToDBConfig())
	defer database.CloseDB(app.db)

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
	), nil, []*river.PeriodicJob{
		river.NewPeriodicJob(
			river.PeriodicInterval(5*time.Minute),
			func() (river.JobArgs, *river.InsertOpts) {
				return jobs.CleanupBlacklistArgs{}, nil
			},
			&river.PeriodicJobOpts{RunOnStart: true},
		),
	})
	if err = app.river.Start(ctx); err != nil {
		errx.Exit(err, "failed to start river client")
	}
	defer app.river.Stop(ctx)

	EnsureKeysExist(ctx, app.db, app.river)

	shutdown := telemetry.InitTracer(ctx, app.cfg.AppName)
	defer telemetry.ShutdownTracer(ctx, shutdown, app.cfg.AppName)

	app.run()
}
