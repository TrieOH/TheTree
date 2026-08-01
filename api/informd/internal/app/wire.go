package app

import (
	"Informd/internal/authz"
	"Informd/internal/handlers"
	"Informd/internal/repos"
	"Informd/internal/services"
	"Informd/internal/sqlc"
	"net/http"
)

// ── Wire types ────────────────────────────────────────────────────────────

type middlewares struct {
	jwt     func(http.Handler) http.Handler
	apiKey  func(http.Handler) http.Handler
	anyAuth func(http.Handler) http.Handler
}

// ── Init functions ────────────────────────────────────────────────────────

func (app *Informd) initRepos(q *sqlc.Queries) *repos.Repos {
	r := repos.New(q)
	authz.Service = authz.New(r.Forms, r.Namespaces)
	return r
}

func (app *Informd) initOperations(r *repos.Repos) *services.Operations {
	return services.NewOperations(r)
}

func (app *Informd) initMiddlewares() middlewares {
	var mw middlewares
	authMW := app.setupAuthMiddlewares()
	mw.jwt = authMW.JWT()
	mw.apiKey = authMW.APIKey()
	mw.anyAuth = authMW.AnyAuth()
	return mw
}

func (app *Informd) initHandlers(ops *services.Operations) *handlers.Server {
	return handlers.NewServer(ops)
}
