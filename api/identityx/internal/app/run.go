package app

import (
	"IdentityX/internal/database/sqlc"
	"lib/database"
	"lib/telemetry"
	"log"
	"net/http"

	"go.opentelemetry.io/otel"
)

func (app *IdentityX) run() {
	q := sqlc.New(app.db)
	loggr := telemetry.Log()
	tx := database.NewPGXTxRunner(app.db, loggr)
	tracer := otel.Tracer(app.cfg.AppName)

	repos := initRepos(q, loggr, tracer)
	queries := initQueries(repos, tx, loggr, tracer)
	commands := initCommands(repos, tx, loggr, tracer)
	handlers := initHandlers(queries, commands)
	middlewares := initMiddlewares(repos, loggr, app.cfg)

	if app.cfg.ProfilePort != "" {
		go servePprof(app.cfg.ProfilePort)
	}

	mux := app.CreateRouter(middlewares, handlers)

	log.Printf("IdentityX listening on :%s", app.cfg.Port)
	log.Fatal(http.ListenAndServe(":"+app.cfg.Port, mux))
}
