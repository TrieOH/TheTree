package app

import (
	"lib/database"
	"lib/telemetry"
	"log"
	"net/http"
	"payssage/internal/database/sqlc"

	"go.opentelemetry.io/otel"
)

func (app *Payssage) run() {
	q := sqlc.New(app.db)
	loggr := telemetry.Log()
	tx := database.NewPGXTxRunner(app.db, loggr)
	tracer := otel.Tracer(app.cfg.AppName)

	repos := initRepos(q, loggr, tracer)

	//river := libriver.NewClient(app.db, libriver.NewWorkers(
	//	libriver.Register[jobs.DeliverWebhookArgs](jobs.NewDeliverWebhookWorker(repos.deliveries)),
	//), nil, nil)
	//if err := river.Start(ctx); err != nil {
	//	telemetry.Log().Fatal("failed to start river client", zap.Error(err))
	//}
	//defer libriver.LogStop(ctx, river, loggr)

	queries := initQueries(repos, app.idxClient, loggr, tx, tracer)
	commands := initCommands(repos, app.idxClient, loggr, tx, tracer)
	middlewares := app.initMiddlewares(loggr, app.cfg)
	handlers := initHandlers(commands, queries)

	if app.cfg.ProfilePort != "" {
		go servePprof(app.cfg.ProfilePort)
	}

	mux := app.CreateRouter(handlers, middlewares)

	log.Printf("Payssage listening on :%s", app.cfg.Port)
	log.Fatal(http.ListenAndServe(":"+app.cfg.Port, mux))
}
