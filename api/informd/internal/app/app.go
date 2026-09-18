package app

import (
	"context"
	"net/http"

	spec "Informd"
	"Informd/internal/config"
	"Informd/internal/sqlc"
	"lib/database"
	"lib/httpserver"

	idx "sdk/identityx"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Informd struct {
	db        *pgxpool.Pool
	idxClient *idx.Client
	cfg       config.Config
}

// Run boots Informd through the Harness: the process sequence (FUN runtime,
// tracer, serving, shutdown ordering) is httpserver.Boot's implementation;
// everything Informd varies on — its IdentityX client, its Postgres adapter
// with the constraint messages, repos, operations, routes — happens in the
// start hook. Storage never crosses Boot's interface: the pool is created
// here and closed in the returned shutdown.
func Run() error {
	cfg := config.LoadConfig()
	app := &Informd{cfg: cfg}

	start := func(ctx context.Context) (http.Handler, func(context.Context) error, error) {
		idxClient, err := idx.Bootstrap(ctx, cfg.ToIdentityXConfig())
		if err != nil {
			return nil, nil, err
		}
		app.idxClient = idxClient

		pool, err := database.SetupDB(cfg.ToDBConfig(), constraintMessages())
		if err != nil {
			return nil, nil, err
		}
		app.db = pool

		tx := database.NewPGXTxRunner(pool)

		repos := app.initRepos(sqlc.New(pool))
		ops := app.initOperations(repos, tx)
		handlers := app.initHandlers(ops)
		primitives := app.initMiddlewares()

		mux := app.CreateRouter(handlers, primitives)
		return mux, func(context.Context) error {
			database.CloseDB(pool)
			return nil
		}, nil
	}

	return httpserver.Boot(httpserver.Config{
		AppName:            cfg.AppName,
		Port:               cfg.Port,
		ProfilePort:        cfg.ProfilePort,
		CorsAllowedOrigins: cfg.CorsAllowedOrigins,
		CorsAllowedHeaders: cfg.CorsAllowedHeaders,
		OpenAPISpec:        spec.OpenAPISpec,
	}, start)
}
