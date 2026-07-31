package app

import (
	"context"
	"lib/database"
	libriver "lib/river"
)

func (app *Univents) run() {
	ctx := context.Background()

	tx := database.NewPGXTxRunner(app.db)
	database.SetDefaultRunner(tx)

	repos := app.initRepos()
	queries := app.initQueries(repos)
	commands := app.initCommands(repos)
	middlewares := app.initMiddlewares()
	handlers := app.initHandlers(queries, commands)

	riverClient, riverUIHandler := app.initRiver(ctx, repos)
	defer libriver.LogStop(ctx, riverClient)

	go servePprof(app.cfg.ProfilePort)
	mux := app.CreateRouter(middlewares, handlers, riverUIHandler)
	app.startServer(mux)
}
