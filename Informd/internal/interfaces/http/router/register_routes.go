package router

import (
	"TrieForms/internal/interfaces/http/middleware"
	"TrieForms/internal/interfaces/http/system"

	"github.com/go-chi/chi/v5"
)

func registerRoutes(r *chi.Mux, deps *HTTPDeps) {
	registerSystemRoutes(r, deps.SystemHandler, deps.AuthMiddleware)
}

func registerSystemRoutes(
	r *chi.Mux,
	h *system.SystemHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Post("/auth/exchange", h.Exchange)
		r.Get("/health", h.Health)
		r.With(authMW.Auth()).
			Get("/protected/health", h.ProtectedHealth)
	})
}
