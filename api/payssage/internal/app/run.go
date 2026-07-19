package app

import (
	"context"
	"lib/database"
	libriver "lib/river"
	"lib/telemetry"
	"log"
	"log/slog"
	"net/http"
	"payssage/internal/database/sqlc"
	webhooksjobs "payssage/internal/features/webhooks/jobs"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"riverqueue.com/riverui"
)

func (app *Payssage) run() {
	ctx := context.Background()

	q := sqlc.New(app.db)
	loggr := telemetry.Log()
	tx := database.NewPGXTxRunner(app.db, loggr)
	tracer := otel.Tracer(app.cfg.AppName)

	repos := initRepos(q, loggr, tracer)
	initProviders(app.cfg.MercadoPagoConfig, repos, loggr, tx, tracer)

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

	queries := initQueries(repos, app.idxClient, loggr, tx, tracer)
	commands := initCommands(riverClient, repos, app.idxClient, loggr, tx, tracer)
	middlewares := app.initMiddlewares(loggr, app.cfg)
	handlers := initHandlers(commands, queries)

	if app.cfg.ProfilePort != "" {
		go servePprof(app.cfg.ProfilePort)
	}
	mux := app.CreateRouter(handlers, middlewares, riverUIHandler)
	log.Printf("Payssage listening on :%s", app.cfg.Port)
	log.Fatal(http.ListenAndServe(":"+app.cfg.Port, mux))
}
