package app

import (
	"context"

	"IdentityX/internal/sqlc"
	"lib/database"
	"lib/httpserver"
	libriver "lib/river"
)

func (app *IdentityX) run() {
	ctx := context.Background()

	q := sqlc.New(app.db)
	tx := database.NewPGXTxRunner(app.db)
	database.SetDefaultRunner(tx)

	riverClient, riverUIHandler := app.initRiver(ctx, q)
	defer libriver.LogStop(ctx, riverClient)
	EnsureKeysExist(ctx, app.db, riverClient)

	repos := app.initRepos(q)
	tokensMgr := app.initTokens(repos)
	ops, authzSvc := app.initOperations(repos, tokensMgr, riverClient)
	handlers := app.initHandlers(ops)
	middlewares := app.initMiddlewares(repos, tokensMgr, authzSvc)

	mux := app.CreateRouter(middlewares, handlers, riverUIHandler)
	httpserver.Start(mux, httpserver.Config{
		AppName:     app.cfg.AppName,
		Port:        app.cfg.Port,
		ProfilePort: app.cfg.ProfilePort,
	})
}
