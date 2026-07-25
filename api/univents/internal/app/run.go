package app

import (
	"lib/database"
	"univents/internal/sqlc"
)

func (app *Univents) run() {
	q := sqlc.New(app.db)
	tx := database.NewPGXTxRunner(app.db)
	database.SetDefaultRunner(tx)

	repos := initRepos(q)
	queries := initQueries(repos)
	commands := initCommands(repos, app.objStorage, app.idxClient)
	middlewares := initMiddlewares()
	handlers := initHandlers(queries, commands)

	go servePprof(app.cfg.ProfilePort)
	mux := app.CreateRouter(middlewares, handlers)
	app.startServer(mux)
}
