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
	// The notifier (lib/go/database) is the store's LISTEN/NOTIFY bridge:
	// the webhook receiver publishes on it (split 4); the SSE relay and WS
	// hub subscribe in split 6. Notify opens its own connection per call,
	// so nothing needs starting here.
	notifier := database.NewNotifier(app.cfg.ToDBConfig().DSN())

	// River must exist before the operations: the webhook receiver cancels
	// the expiry job on approve via the client (best-effort; split 7
	// checkout schedules the job).
	riverClient, riverUIHandler := app.initRiver(ctx, repos, notifier, tx)
	defer libriver.LogStop(ctx, riverClient)

	ops := app.initOperations(repos, notifier, riverClient, tx)

	middlewares := app.initMiddlewares()
	handlers := app.initHandlers(ops)

	mux := app.CreateRouter(middlewares, handlers, riverUIHandler)
	httpserver.Start(mux, httpserver.Config{
		AppName:     app.cfg.AppName,
		Port:        app.cfg.Port,
		ProfilePort: app.cfg.ProfilePort,
	})
}
