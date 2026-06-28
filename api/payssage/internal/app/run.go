package app

import (
	"context"
	"lib/database"
	libriver "lib/river"
	"lib/telemetry"
	"log"
	"net/http"
	"payssage/internal/database/sqlc"
	"payssage/internal/jobs"
	"payssage/ports"

	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type paymentProviders struct {
	oauth    map[string]ports.OAuthProvider
	payments map[string]ports.PaymentAbstractionLayer
}

func (app *Payssage) run() {
	ctx := context.Background()
	q := sqlc.New(app.db)
	loggr := telemetry.Log()
	tx := database.NewPGXTxRunner(app.db, loggr)
	tracer := otel.Tracer(app.cfg.AppName)
	paymentProviders := setupPaymentProviders(app.cfg)

	repos := initRepos(q, loggr, tracer)

	river := libriver.NewClient(app.db, libriver.NewWorkers(
		libriver.Register[jobs.DeliverWebhookArgs](jobs.NewDeliverWebhookWorker(repos.deliveries)),
	), nil, nil)
	if err := river.Start(ctx); err != nil {
		telemetry.Log().Fatal("failed to start river client", zap.Error(err))
	}
	defer libriver.LogStop(ctx, river, loggr)

	queries := initQueries(repos, loggr, tx, tracer)
	commands := initCommands(repos, river, loggr, tx, tracer)
	handlers := initHandlers(commands, queries)

	if app.cfg.ProfilePort != "" {
		go servePprof(app.cfg.ProfilePort)
	}

	mux := app.CreateRouter(handlers)

	log.Printf("Payssage listening on :%s", app.cfg.Port)
	log.Fatal(http.ListenAndServe(":"+app.cfg.Port, mux))
}
