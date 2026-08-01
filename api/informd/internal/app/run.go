package app

import (
	"Informd/internal/sqlc"
	"lib/database"
	"lib/httpserver"
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

	mux := app.CreateRouter(handlers, middlewares)
	httpserver.Start(mux, httpserver.Config{
		AppName:     app.cfg.AppName,
		Port:        app.cfg.Port,
		ProfilePort: app.cfg.ProfilePort,
	})
}
