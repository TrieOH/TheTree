package app

import (
	"context"

	"IdentityX/internal/sqlc"
	"lib/database"
	"lib/errx"
	"lib/httpserver"
	libriver "lib/river"
)

func (app *IdentityX) run() {
	ctx := context.Background()

	q := sqlc.New(app.db)
	tx := database.NewPGXTxRunner(app.db)
	database.SetDefaultRunner(tx)

	repos := app.initRepos(q)
	actionTokenMgr := app.initActionTokens(repos)
	keysMgr := app.initKeys(repos)

	// Provision every scope's keys before the router accepts traffic: the
	// Key-lifecycle module creates what is missing, rotates expired or
	// legacy no-expiry keys, and sweeps retiring keys. The periodic
	// RotateKeysWorker keeps them fresh afterwards.
	err := keysMgr.EnsureAll(ctx)
	if err != nil {
		errx.Exit(err, "failed to ensure crypto keys")
	}

	riverClient, riverUIHandler := app.initRiver(ctx, q, actionTokenMgr, keysMgr)
	defer libriver.LogStop(ctx, riverClient)

	tokensMgr := app.initTokens(repos)
	ops, authzSvc := app.initOperations(repos, tokensMgr, actionTokenMgr, keysMgr, riverClient)
	handlers := app.initHandlers(ops)
	middlewares := app.initMiddlewares(ops, tokensMgr, authzSvc)

	mux := app.CreateRouter(middlewares, handlers, riverUIHandler)
	httpserver.Start(mux, httpserver.Config{
		AppName:     app.cfg.AppName,
		Port:        app.cfg.Port,
		ProfilePort: app.cfg.ProfilePort,
	})
}
