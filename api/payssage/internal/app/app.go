package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"lib/database"
	"lib/httpserver"
	libriver "lib/river"
	spec "payssage"
	"payssage/internal/config"
	webhooksjobs "payssage/internal/services/webhooks/jobs"
	"payssage/internal/sqlc"

	idx "sdk/identityx"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"resty.dev/v3"
	"riverqueue.com/riverui"
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
// its Postgres adapter with the constraint messages, its webhook delivery
// workers and the river UI — happens in the start hook. Storage never
// crosses Boot's interface: the pool is created here and closed in the
// returned shutdown.
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

		pool, err := database.SetupDB(cfg.ToDBConfig(), constraintMessages())
		if err != nil {
			return nil, nil, err
		}
		app.db = pool

		tx := database.NewPGXTxRunner(pool)

		repos := app.initRepos(sqlc.New(pool))
		app.initProviders(repos)

		riverClient := libriver.NewClient(pool, libriver.NewWorkers(
			libriver.Register[webhooksjobs.DeliverWebhookArgs](webhooksjobs.NewDeliverWebhookWorker(
				repos.WebhookDeliveries, repos.WebhookEvents, repos.WebhookEndpoints, app.httpClient,
			)),
		), nil, nil)
		err = riverClient.Start(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("start river client: %w", err)
		}

		riverUIHandler, err := riverui.NewHandler(&riverui.HandlerOpts{
			DevMode:                  false,
			Endpoints:                riverui.NewEndpoints[pgx.Tx](riverClient, nil),
			Logger:                   slog.Default(),
			Prefix:                   "/riverui",
			JobListHideArgsByDefault: true,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("create river ui handler: %w", err)
		}
		err = riverUIHandler.Start(ctx)
		if err != nil {
			return nil, nil, fmt.Errorf("start river ui handler: %w", err)
		}

		ops := app.initOperations(riverClient, repos, tx)
		handlers := app.initHandlers(ops)
		primitives := app.initMiddlewares()

		mux := app.CreateRouter(primitives, handlers, riverUIHandler)
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
		CorsAllowedOrigins: cfg.CorsAllowedOrigins,
		CorsAllowedHeaders: cfg.CorsAllowedHeaders,
		OpenAPISpec:        spec.OpenAPISpec,
	}, start)
}
