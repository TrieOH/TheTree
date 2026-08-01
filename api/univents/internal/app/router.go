package app

import (
	"net/http"
	"univents/internal/features/badges"
	"univents/internal/features/certifications"
	"univents/internal/features/editions"
	"univents/internal/features/events"
	"univents/internal/features/products"
	"univents/internal/features/programs"
	"univents/internal/features/signatures"
	"univents/internal/features/ticket_types"

	"lib/httpserver"

	"github.com/go-chi/chi/v5"
	"riverqueue.com/riverui"
)

func (app *Univents) CreateRouter(middlewares middlewares, handlers handlers, riverUIHandler *riverui.Handler) http.Handler {
	return httpserver.NewRouter(httpserver.Config{
		AppName:         app.cfg.AppName,
		SkipLogPrefixes: []string{"/admin/asynq"},
		Routes: func(r *chi.Mux) {
			events.RegisterRoutes(r, handlers.events, middlewares.jwt)
			editions.RegisterRoutes(r, handlers.editions, middlewares.jwt)
			ticket_types.RegisterRoutes(r, handlers.ticketTypes, middlewares.jwt)
			products.RegisterRoutes(r, handlers.products, middlewares.jwt)
			programs.RegisterRoutes(r, handlers.programs, middlewares.jwt)
			badges.RegisterRoutes(r, handlers.badges, middlewares.jwt)
			signatures.RegisterRoutes(r, handlers.signatures, middlewares.jwt)
			certifications.RegisterRoutes(r, handlers.certs, middlewares.jwt)

			r.Group(func(r chi.Router) {
				r.Mount("/riverui", riverUIHandler)
			})
		},
	})
}
