package app

import (
	"IdentityX/internal/sqlc"
	"lib/database"
	"lib/httpserver"
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

	mux := app.CreateRouter(middlewares, handlers)
	httpserver.Start(mux, httpserver.Config{
		AppName:     app.cfg.AppName,
		Port:        app.cfg.Port,
		ProfilePort: app.cfg.ProfilePort,
	})
}
