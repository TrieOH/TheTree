package app

import (
	"context"
	"lib/database"
	libriver "lib/river"
	"lib/telemetry"
	"log/slog"
	webhooksjobs "payssage/internal/features/webhooks/jobs"
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
			repos.deliveries, repos.events, repos.endpoints,
		)),
	), nil, nil)
	if err := riverClient.Start(ctx); err != nil {
		loggr.Fatal("failed to start river client", zap.Error(err))
	}
	defer libriver.LogStop(ctx, riverClient, loggr)

	riverUIHandler, err := riverui.NewHandler(&riverui.HandlerOpts{
		DevMode:   false,
		Endpoints: riverui.NewEndpoints[pgx.Tx](riverClient, nil),
		Logger:    slog.Default(),
		Prefix:    "/riverui",
	})
	if err != nil {
		loggr.Fatal("failed to create river ui handler", zap.Error(err))
	}
	if err := riverUIHandler.Start(ctx); err != nil {
		loggr.Fatal("failed to start river ui handler", zap.Error(err))
	}

	queries := app.initQueries(repos)
	commands := app.initCommands(riverClient, repos)
	middlewares := app.initMiddlewares()
	handlers := app.initHandlers(commands, queries)

	if app.cfg.ProfilePort != "" {
		go servePprof(app.cfg.ProfilePort)
	}
	mux := app.CreateRouter(handlers, middlewares, riverUIHandler)
	app.startServer(mux)
}
