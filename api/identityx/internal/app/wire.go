package app

import (
	"IdentityX/internal/authz"
	"IdentityX/internal/handlers"
	"IdentityX/internal/repos"
	"IdentityX/internal/services"
	"IdentityX/internal/sqlc"
	"net/http"
)

// ── Wire types ────────────────────────────────────────────────────────────

type middlewares struct {
	jwtAuth           func(http.Handler) http.Handler
	apiKeyAuth        func(http.Handler) http.Handler
	anyAuth           func(http.Handler) http.Handler
	clientOnly        func(http.Handler) http.Handler
	projectClientOnly func(http.Handler) http.Handler
}

// ── Init functions ────────────────────────────────────────────────────────

func (app *IdentityX) initRepos(q *sqlc.Queries) *repos.Repos {
	return repos.New(q)
}

func (app *IdentityX) initOperations(r *repos.Repos) *services.Operations {
	authzSvc := authz.New(r.Organizations, r.Projects)
	return services.NewOperations(r, authzSvc, app.cfg.HmacSecret)
}

func (app *IdentityX) initMiddlewares(r *repos.Repos) middlewares {
	var mw middlewares
	authMW := app.SetupAuthMiddlewares(r.CryptoKeys, r.APIKeys, r.Actors, r.Capabilities, r.Blacklist)
	mw.jwtAuth = authMW.JWT()
	mw.apiKeyAuth = authMW.APIKey()
	mw.anyAuth = authMW.AnyAuth()
	mw.clientOnly = ClientOnly()
	mw.projectClientOnly = ProjectClientOnly()
	return mw
}

func (app *IdentityX) initHandlers(ops *services.Operations) *handlers.Server {
	return handlers.NewServer(ops)
}
