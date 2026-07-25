package app

import (
	"Informd/internal/sqlc"
	"lib/database"
)

func (app *Informd) run() {
	q := sqlc.New(app.db)
	tx := database.NewPGXTxRunner(app.db)
	database.SetDefaultRunner(tx)

	repos := app.initRepos(q)
	queries := app.initQueries(repos)
	commands := app.initCommands(repos)
	handlers := app.initHandlers(commands, queries)
	middlewares := app.initMiddlewares()

	if app.cfg.ProfilePort != "" {
		go servePprof(app.cfg.ProfilePort)
	}

	mux := app.CreateRouter(handlers, middlewares)
	app.startServer(mux)
}
