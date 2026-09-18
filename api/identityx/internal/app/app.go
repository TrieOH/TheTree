package app

import (
	"context"
	"fmt"
	"net/http"

	spec "IdentityX"
	"IdentityX/internal/config"
	"IdentityX/internal/setup"
	"IdentityX/internal/sqlc"
	"lib/database"
	"lib/email"
	"lib/httpserver"
	libriver "lib/river"

	"github.com/jackc/pgx/v5/pgxpool"
)

type IdentityX struct {
	db          *pgxpool.Pool
	emailClient *email.Client

	cfg config.Config
}

// Run boots IdentityX through the Harness: the process sequence (FUN
// runtime, tracer, serving, shutdown ordering) is httpserver.Boot's
// implementation; everything IdentityX varies on — its Postgres adapter
// with the constraint messages, the setup-state flag, the Key-lifecycle
// provisioning, river workers and jobs — happens in the start hook.
// Storage never crosses Boot's interface: the pool is created here and
// closed in the returned shutdown.
func Run() error {
	cfg := config.LoadConfig()
	app := &IdentityX{cfg: cfg}

	start := func(ctx context.Context) (http.Handler, func(context.Context) error, error) {
		pool, err := database.SetupDB(cfg.ToDBConfig(), constraintMessages())
		if err != nil {
			return nil, nil, err
		}
		app.db = pool
		app.emailClient = email.NewClient(cfg.ToEmailConfig())

		q := sqlc.New(pool)
		has, err := q.HasAnyActor(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("check setup state: %w", err)
		}
		if has {
			setup.MarkComplete()
		}

		tx := database.NewPGXTxRunner(pool)

		repos := app.initRepos(q)
		actionTokenMgr := app.initActionTokens(repos)
		keysMgr := app.initKeys(repos)

		// Provision every scope's keys before the router accepts traffic:
		// the Key-lifecycle module creates what is missing, rotates expired
		// or legacy no-expiry keys, and sweeps retiring keys. The periodic
		// RotateKeysWorker keeps them fresh afterwards.
		err = keysMgr.EnsureAll(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("ensure crypto keys: %w", err)
		}

		riverClient, riverUI, err := app.initRiver(ctx, q, actionTokenMgr, keysMgr)
		if err != nil {
			return nil, nil, err
		}

		tokensMgr := app.initTokens(repos)
		ops, authzSvc := app.initOperations(repos, tokensMgr, actionTokenMgr, keysMgr, riverClient, tx)
		handlers := app.initHandlers(ops)
		primitives := app.initMiddlewares(ops, tokensMgr, authzSvc)

		mux := app.CreateRouter(primitives, handlers, riverUI)
		return mux, func(ctx context.Context) error {
			libriver.LogStop(ctx, riverClient)
			database.CloseDB(pool)
			return nil
		}, nil
	}

	return httpserver.Boot(httpserver.Config{
		AppName:            cfg.AppName,
		Port:               cfg.Port,
		ProfilePort:        cfg.ProfilePort,
		CorsAllowedOrigins: cfg.AllowedOrigins,
		CorsAllowedHeaders: cfg.AllowedHeaders,
		OpenAPISpec:        spec.OpenAPISpec,
	}, start)
}
