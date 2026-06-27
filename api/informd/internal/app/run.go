package app

import (
	"Informd/internal/database/sqlc"
	"lib/database"
	"lib/telemetry"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
)

func (app *Informd) run() {
	q := sqlc.New(app.db)
	loggr := telemetry.Log()
	tx := database.NewPGXTxRunner(app.db, loggr)
	tracer := otel.Tracer(app.cfg.AppName)

	repos := initRepos(q, loggr, tracer)
	queries := initQueries(repos, loggr, tracer, tx)
	commands := initCommands(repos, loggr, tracer, tx)
	handlers := initHandlers(commands, queries)
	middlewares := initMiddlewares(loggr)

	if app.cfg.ProfilePort != "" {
		go servePprof(app.cfg.ProfilePort)
	}

	mux := app.CreateRouter(handlers, middlewares)

	log.Printf("Informd listening on :%s", app.cfg.Port)
	log.Fatal(http.ListenAndServe(":"+app.cfg.Port, mux))
}
