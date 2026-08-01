package app

import (
	"IdentityX/internal/features/actors"
	"IdentityX/internal/features/api_keys"
	"IdentityX/internal/features/authn"
	"IdentityX/internal/features/capabilities"
	"IdentityX/internal/features/organizations"
	"IdentityX/internal/features/profile_schemas"
	"IdentityX/internal/features/profiles"
	"IdentityX/internal/features/projects"
	"net/http"

	"lib/httpserver"

	"github.com/go-chi/chi/v5"
)

func (app *IdentityX) CreateRouter(middlewares middlewares, handlers handlers) http.Handler {
	return httpserver.NewRouter(httpserver.Config{
		AppName: app.cfg.AppName,
		Routes: func(r *chi.Mux) {
			actors.RegisterRoutes(r, handlers.Actors, middlewares.anyAuth, middlewares.clientOnly)
			api_keys.RegisterRoutes(r, handlers.APIKeys, middlewares.anyAuth, middlewares.clientOnly)
			authn.RegisterRoutes(r, handlers.Authn, middlewares.jwtAuth, middlewares.anyAuth)
			organizations.RegisterRoutes(r, handlers.Orgs, middlewares.jwtAuth, middlewares.clientOnly)
			projects.RegisterRoutes(r, handlers.Projects, middlewares.anyAuth, middlewares.clientOnly)
			capabilities.RegisterRoutes(r, handlers.Capabilities, middlewares.jwtAuth, middlewares.anyAuth, middlewares.clientOnly)
			profiles.RegisterRoutes(r, handlers.Profiles, middlewares.jwtAuth, middlewares.clientOnly)
			profile_schemas.RegisterRoutes(r, handlers.ProfileSchemas, middlewares.jwtAuth, middlewares.clientOnly)
		},
	})
}
