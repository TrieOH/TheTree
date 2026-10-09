package app

import (
	"context"
	"fmt"
	"net/http"

	spec "IdentityX"
	"IdentityX/internal/config"
	"IdentityX/internal/jobs"
	"IdentityX/internal/setup"
	"IdentityX/internal/sqlc"
	"lib/database"
	"lib/email"
	"lib/httpserver"
	"lib/jobs/lambdaevents"

	"github.com/jackc/pgx/v5/pgxpool"
)

type IdentityX struct {
	db    *pgxpool.Pool
	email email.Sender

	cfg config.Config
}

// Run boots IdentityX through the Harness: the process sequence (FUN
// runtime, telemetry, serving, shutdown ordering) is httpserver.Boot's
// implementation; everything IdentityX varies on — its Postgres adapter
// with the constraint messages, the setup-state probe, the Key-lifecycle
// provisioning, the email transport and the background-job carrier —
// happens in the start hook. Storage never crosses Boot's interface: the
// pool is created here and closed in the returned shutdown.
//
// The runtime mode decides the carrier: hosted runs jobs in-process
// (workers plus tickers); managed sends them to a queue drained by the
// worker role, whose periodic jobs a scheduler fires (see `identityx
// schedules`). The schema is never migrated here — `just identityx-goose`
// owns it.
func Run() error {
	cfg := config.LoadConfig()
	app := &IdentityX{cfg: cfg}

	start := func(ctx context.Context) (http.Handler, func(context.Context) error, error) {
		pool, err := database.SetupDBWithoutMigrations(cfg.ToDBConfig(), constraintMessages())
		if err != nil {
			return nil, nil, err
		}
		app.db = pool

		app.email, err = app.initEmail(ctx)
		if err != nil {
			database.CloseDB(pool)
			return nil, nil, err
		}

		q := sqlc.New(pool)
		setup.UseProbe(q.HasAnyActor)

		tx := database.NewPGXTxRunner(pool)

		repos := app.initRepos(q)
		actionTokenMgr := app.initActionTokens(repos)
		keysMgr := app.initKeys(repos)

		// Provision every scope's keys before the router accepts traffic:
		// the Key-lifecycle module creates what is missing, rotates expired
		// or legacy no-expiry keys, and sweeps retiring keys. The periodic
		// RotateKeys job keeps them fresh afterwards.
		err = keysMgr.EnsureAll(ctx)
		if err != nil {
			database.CloseDB(pool)
			return nil, nil, fmt.Errorf("ensure crypto keys: %w", err)
		}

		registry := jobs.Registry(jobs.Deps{
			Queries:      q,
			ActionTokens: actionTokenMgr,
			Keys:         keysMgr,
			Email:        app.email,
		})
		enqueuer, stopJobs, err := app.initJobs(ctx, registry)
		if err != nil {
			database.CloseDB(pool)
			return nil, nil, err
		}

		tokensMgr := app.initTokens(repos)
		ops, authzSvc := app.initOperations(repos, tokensMgr, actionTokenMgr, keysMgr, enqueuer, tx)
		handlers := app.initHandlers(ops)
		primitives := app.initMiddlewares(ops, tokensMgr, authzSvc)

		var events http.Handler
		if cfg.Worker() {
			events = lambdaevents.Handler(registry)
		}

		mux := app.CreateRouter(primitives, handlers, events)
		return mux, func(ctx context.Context) error {
			err := stopJobs(ctx)
			database.CloseDB(pool)
			return err
		}, nil
	}

	profilePort := cfg.ProfilePort
	if cfg.Managed() {
		// A function has no reachable side port.
		profilePort = ""
	}
	return httpserver.Boot(httpserver.Config{
		AppName:            cfg.AppName,
		Port:               cfg.Port,
		ProfilePort:        profilePort,
		CorsAllowedOrigins: cfg.AllowedOrigins,
		CorsAllowedHeaders: cfg.AllowedHeaders,
		OpenAPISpec:        spec.OpenAPISpec,
		RateLimit:          cfg.ToRateLimit(),
	}, start)
}
