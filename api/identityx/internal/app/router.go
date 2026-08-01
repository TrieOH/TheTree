package app

import (
	actorsHandlers "IdentityX/internal/features/actors/handlers"
	apiKeysHandlers "IdentityX/internal/features/api_keys/handlers"
	authnHandlers "IdentityX/internal/features/authn/handlers"
	capabilitiesHandlers "IdentityX/internal/features/capabilities/handlers"
	orgsHandlers "IdentityX/internal/features/organizations/handlers"
	profileSchemasHandlers "IdentityX/internal/features/profile_schemas/handlers"
	profilesHandlers "IdentityX/internal/features/profiles/handlers"
	projectsHandlers "IdentityX/internal/features/projects/handlers"
	"net/http"

	spec "IdentityX"

	"lib/httpserver"

	"github.com/go-chi/chi/v5"
)

func (app *IdentityX) CreateRouter(middlewares middlewares, handlers handlers) http.Handler {
	return httpserver.NewRouter(httpserver.Config{
		AppName:     app.cfg.AppName,
		OpenAPISpec: spec.OpenAPISpec,
		Routes: func(r *chi.Mux) {
			registerRoutes(r, middlewares, handlers)
		},
	})
}

// registerRoutes wires every feature's routes onto r. Kept package-level so
// the router-parity test can walk the same registration the app serves.
func registerRoutes(r *chi.Mux, middlewares middlewares, handlers handlers) {
	actorsHandlers.RegisterRoutes(r, handlers.Actors, middlewares.anyAuth, middlewares.clientOnly)
	apiKeysHandlers.RegisterRoutes(r, handlers.APIKeys, middlewares.anyAuth, middlewares.clientOnly)
	authnHandlers.RegisterRoutes(r, handlers.Authn, middlewares.jwtAuth, middlewares.anyAuth)
	orgsHandlers.RegisterRoutes(r, handlers.Orgs, middlewares.jwtAuth, middlewares.clientOnly)
	projectsHandlers.RegisterRoutes(r, handlers.Projects, middlewares.anyAuth, middlewares.clientOnly)
	capabilitiesHandlers.RegisterRoutes(r, handlers.Capabilities, middlewares.jwtAuth, middlewares.anyAuth, middlewares.clientOnly)
	profilesHandlers.RegisterRoutes(r, handlers.Profiles, middlewares.jwtAuth, middlewares.clientOnly)
	profileSchemasHandlers.RegisterRoutes(r, handlers.ProfileSchemas, middlewares.jwtAuth, middlewares.clientOnly)
}
