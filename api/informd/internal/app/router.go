package app

import (
	fieldsHandlers "Informd/internal/features/fields/handlers"
	formsHandlers "Informd/internal/features/forms/handlers"
	namespacesHandlers "Informd/internal/features/namespaces/handlers"
	responsesHandlers "Informd/internal/features/responses/handlers"
	stepsHandlers "Informd/internal/features/steps/handlers"
	"net/http"

	spec "Informd"

	"lib/httpserver"

	"github.com/go-chi/chi/v5"
)

func (app *Informd) CreateRouter(handlers handlers, middlewares middlewares) http.Handler {
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
	namespacesHandlers.RegisterRoutes(r, handlers.namespaces, middlewares.jwt)
	formsHandlers.RegisterRoutes(r, handlers.forms, middlewares.anyAuth)
	stepsHandlers.RegisterRoutes(r, handlers.steps, middlewares.anyAuth)
	fieldsHandlers.RegisterRoutes(r, handlers.fields, middlewares.anyAuth)
	responsesHandlers.RegisterRoutes(r, handlers.responses)
}
