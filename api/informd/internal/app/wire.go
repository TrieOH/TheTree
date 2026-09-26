package app

import (
	"Informd/internal/authz"
	"Informd/internal/handlers"
	"Informd/internal/repos"
	"Informd/internal/services"
	"Informd/internal/sqlc"
	"lib/database"

	libauthz "lib/authz"
)

// ── Init functions ────────────────────────────────────────────────────────

func (app *Informd) initRepos(q *sqlc.Queries) *repos.Repos {
	return repos.New(q)
}

func (app *Informd) initOperations(r *repos.Repos, tx database.TxRunner) *services.Operations {
	authzSvc := authz.New(r.Forms, r.Namespaces)
	return services.NewOperations(r, authzSvc, tx)
}

// initMiddlewares builds the auth primitives the spec-derived chains
// resolve against; construction stays per-backend, everything downstream
// (chain derivation, dispatch, fail-closed) is the Access-check and
// Harness modules' implementation.
func (app *Informd) initMiddlewares() libauthz.Primitives {
	authMW := app.setupAuthMiddlewares()
	return libauthz.Primitives{
		JWT:    authMW.JWT(),
		APIKey: authMW.APIKey(),
		Any:    authMW.AnyAuth(),
	}
}

func (app *Informd) initHandlers(ops *services.Operations) *handlers.Server {
	return handlers.NewServer(ops)
}
