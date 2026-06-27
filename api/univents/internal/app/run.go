package app

import (
	"lib/database"
	"lib/telemetry"
	"log"
	"net/http"
	"univents/internal/platform/database/sqlc"

	"go.opentelemetry.io/otel"
)

func (app *Univents) run() {
	q := sqlc.New(app.db)
	loggr := telemetry.Log()
	tx := database.NewPGXTxRunner(app.db, loggr)
	tracer := otel.Tracer(app.cfg.AppName)

	//rt.wsRegistry = sockets.New()

	repos := initRepos(q, loggr, tracer)
	queries := initQueries(repos, tx, loggr, tracer)
	commands := initCommands(repos, tx, loggr, tracer)
	middlewares := initMiddlewares(loggr)
	handlers := initHandlers(queries, commands)

	if app.cfg.ProfilePort != "" {
		go servePprof(app.cfg.ProfilePort)
	}

	mux := app.CreateRouter(middlewares, handlers)

	log.Printf("IdentityX listening on :%s", app.cfg.Port)
	log.Fatal(http.ListenAndServe(":"+app.cfg.Port, mux))
}
