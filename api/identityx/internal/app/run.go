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
	ops := app.initOperations(repos)
	handlers := app.initHandlers(ops)
	middlewares := app.initMiddlewares(repos)

	mux := app.CreateRouter(middlewares, handlers)
	httpserver.Start(mux, httpserver.Config{
		AppName:     app.cfg.AppName,
		Port:        app.cfg.Port,
		ProfilePort: app.cfg.ProfilePort,
	})
}
