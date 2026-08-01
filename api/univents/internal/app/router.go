package app

import (
	"net/http"
	badgesHandlers "univents/internal/features/badges/handlers"
	certificationsHandlers "univents/internal/features/certifications/handlers"
	editionsHandlers "univents/internal/features/editions/handlers"
	eventsHandlers "univents/internal/features/events/handlers"
	productsHandlers "univents/internal/features/products/handlers"
	programsHandlers "univents/internal/features/programs/handlers"
	signaturesHandlers "univents/internal/features/signatures/handlers"
	ticketTypesHandlers "univents/internal/features/ticket_types/handlers"

	"lib/httpserver"
	spec "univents"

	"github.com/go-chi/chi/v5"
	"riverqueue.com/riverui"
)

func (app *Univents) CreateRouter(middlewares middlewares, handlers handlers, riverUIHandler *riverui.Handler) http.Handler {
	return httpserver.NewRouter(httpserver.Config{
		AppName:         app.cfg.AppName,
		OpenAPISpec:     spec.OpenAPISpec,
		SkipLogPrefixes: []string{"/admin/asynq"},
		Routes: func(r *chi.Mux) {
			registerRoutes(r, middlewares, handlers, riverUIHandler)
		},
	})
}

// registerRoutes wires every feature's routes onto r. Kept package-level so
// the router-parity test can walk the same registration the app serves.
func registerRoutes(r *chi.Mux, middlewares middlewares, handlers handlers, riverUIHandler *riverui.Handler) {
	eventsHandlers.RegisterRoutes(r, handlers.events, middlewares.jwt)
	editionsHandlers.RegisterRoutes(r, handlers.editions, middlewares.jwt)
	ticketTypesHandlers.RegisterRoutes(r, handlers.ticketTypes, middlewares.jwt)
	productsHandlers.RegisterRoutes(r, handlers.products, middlewares.jwt)
	programsHandlers.RegisterRoutes(r, handlers.programs, middlewares.jwt)
	badgesHandlers.RegisterRoutes(r, handlers.badges, middlewares.jwt)
	signaturesHandlers.RegisterRoutes(r, handlers.signatures, middlewares.jwt)
	certificationsHandlers.RegisterRoutes(r, handlers.certs, middlewares.jwt)

	r.Group(func(r chi.Router) {
		r.Mount("/riverui", riverUIHandler)
	})
}
