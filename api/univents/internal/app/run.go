package app

import (
	"context"
	"lib/database"
	"lib/httpserver"
	libriver "lib/river"
)

func (app *Univents) run() {
	ctx := context.Background()

	tx := database.NewPGXTxRunner(app.db)
	database.SetDefaultRunner(tx)

	repos := app.initRepos()
	ops := app.initOperations(repos)

	middlewares := app.initMiddlewares()
	handlers := app.initHandlers(ops)

	riverClient, riverUIHandler := app.initRiver(ctx, repos)
	defer libriver.LogStop(ctx, riverClient)

	mux := app.CreateRouter(middlewares, handlers, riverUIHandler)
	httpserver.Start(mux, httpserver.Config{
		AppName:     app.cfg.AppName,
		Port:        app.cfg.Port,
		ProfilePort: app.cfg.ProfilePort,
	})
}
