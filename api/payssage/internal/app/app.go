package app

import (
	"context"
	"net/http"
	"time"

	"lib/database"
	"lib/httpserver"
	spec "payssage"
	"payssage/internal/config"
	"payssage/internal/sqlc"

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

// Run boots Payssage through the Harness: the process sequence (FUN
// runtime, tracer, serving, shutdown ordering) is httpserver.Boot's
// implementation; everything Payssage varies on — its IdentityX client,
// its Postgres adapter with the constraint messages and the background-job
// carrier — happens in the start hook. Storage never crosses Boot's
// interface: the pool is created here and closed in the returned shutdown.
//
// Background jobs run in-process locally, or through an external queue and
// scheduler when JOBS_URL is set (see initJobs). The schema is never
// migrated here — `just payssage-goose` owns it.
func Run() error {
	cfg := config.LoadConfig()
	app := &Payssage{cfg: cfg}

	start := func(ctx context.Context) (http.Handler, func(context.Context) error, error) {
		idxClient, err := idx.Bootstrap(ctx, cfg.ToIdentityXConfig())
		if err != nil {
			return nil, nil, err
		}
		app.idxClient = idxClient
		app.httpClient = resty.New().SetTimeout(15 * time.Second)

		pool, err := database.SetupDBWithoutMigrations(cfg.ToDBConfig(), constraintMessages())
		if err != nil {
			return nil, nil, err
		}
		app.db = pool

		tx := database.NewPGXTxRunner(pool)

		repos := app.initRepos(sqlc.New(pool))
		app.initProviders(repos)

		enqueuer, events, stopJobs, err := app.initJobs(repos)
		if err != nil {
			database.CloseDB(pool)
			return nil, nil, err
		}

		ops := app.initOperations(enqueuer, repos, tx)
		handlers := app.initHandlers(ops)
		primitives := app.initMiddlewares()

		mux := app.CreateRouter(primitives, handlers, events)
		return mux, func(ctx context.Context) error {
			err := stopJobs(ctx)
			database.CloseDB(pool)
			return err
		}, nil
	}

	return httpserver.Boot(httpserver.Config{
		AppName:            cfg.AppName,
		Port:               cfg.Port,
		ProfilePort:        cfg.ProfilePort,
		CorsAllowedOrigins: cfg.AllowedOrigins,
		CorsAllowedHeaders: cfg.AllowedHeaders,
		OpenAPISpec:        spec.OpenAPISpec,
		RateLimit:          cfg.ToRateLimit(),
	}, start)
}
