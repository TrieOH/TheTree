package app

import (
	"IdentityX/internal/sqlc"
	"lib/database"
)

func (app *IdentityX) run() {
	q := sqlc.New(app.db)
	tx := database.NewPGXTxRunner(app.db)
	database.SetDefaultRunner(tx)

	repos := app.initRepos(q)
	queries := app.initQueries(repos)
	commands := app.initCommands(repos)
	handlers := app.initHandlers(queries, commands)
	middlewares := app.initMiddlewares(repos)

	if app.cfg.ProfilePort != "" {
		go servePprof(app.cfg.ProfilePort)
	}

	mux := app.CreateRouter(middlewares, handlers)
	app.startServer(mux)
}
