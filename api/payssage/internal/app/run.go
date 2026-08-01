package app

import (
	"context"
	"lib/database"
	"lib/httpserver"
	libriver "lib/river"
	"lib/telemetry"
	"log/slog"
	webhooksjobs "payssage/internal/services/webhooks/jobs"
	"payssage/internal/sqlc"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
	"riverqueue.com/riverui"
)

func (app *Payssage) run() {
	ctx := context.Background()

	q := sqlc.New(app.db)
	loggr := telemetry.Log()
	tx := database.NewPGXTxRunner(app.db)
	database.SetDefaultRunner(tx)

	repos := app.initRepos(q)
	app.initProviders(repos)

	riverClient := libriver.NewClient(app.db, libriver.NewWorkers(
		libriver.Register[webhooksjobs.DeliverWebhookArgs](webhooksjobs.NewDeliverWebhookWorker(
			repos.WebhookDeliveries, repos.WebhookEvents, repos.WebhookEndpoints, app.httpClient,
		)),
	), nil, nil)
	err := riverClient.Start(ctx)
	if err != nil {
		loggr.Fatal("failed to start river client", zap.Error(err))
	}
	defer libriver.LogStop(ctx, riverClient)

	riverUIHandler, err := riverui.NewHandler(&riverui.HandlerOpts{
		DevMode:   false,
		Endpoints: riverui.NewEndpoints[pgx.Tx](riverClient, nil),
		Logger:    slog.Default(),
		Prefix:    "/riverui",
	})
	if err != nil {
		loggr.Fatal("failed to create river ui handler", zap.Error(err))
	}
	err = riverUIHandler.Start(ctx)
	if err != nil {
		loggr.Fatal("failed to start river ui handler", zap.Error(err))
	}

	ops := app.initOperations(riverClient, repos)

	middlewares := app.initMiddlewares()
	handlers := app.initHandlers(ops)

	mux := app.CreateRouter(middlewares, handlers, riverUIHandler)
	httpserver.Start(mux, httpserver.Config{
		AppName:     app.cfg.AppName,
		Port:        app.cfg.Port,
		ProfilePort: app.cfg.ProfilePort,
	})
}
