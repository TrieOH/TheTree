package app

import (
	"Informd/internal/features/fields"
	"Informd/internal/features/forms"
	"Informd/internal/features/namespaces"
	"Informd/internal/features/responses"
	"Informd/internal/features/steps"
	"net/http"

	"lib/httpserver"

	"github.com/go-chi/chi/v5"
)

func (app *Informd) CreateRouter(handlers handlers, middlewares middlewares) http.Handler {
	return httpserver.NewRouter(httpserver.Config{
		AppName: app.cfg.AppName,
		Routes: func(r *chi.Mux) {
			namespaces.RegisterRoutes(r, handlers.namespaces, middlewares.jwt)
			forms.RegisterRoutes(r, handlers.forms, middlewares.anyAuth)
			steps.RegisterRoutes(r, handlers.steps, middlewares.anyAuth)
			fields.RegisterRoutes(r, handlers.fields, middlewares.anyAuth)
			responses.RegisterRoutes(r, handlers.responses)
		},
	})
}
