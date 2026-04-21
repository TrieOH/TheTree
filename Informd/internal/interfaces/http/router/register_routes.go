package router

import (
	"TrieForms/internal/features/forms"
	"TrieForms/internal/features/keys"
	"TrieForms/internal/features/projects"
	"TrieForms/internal/interfaces/http/middleware"
	"TrieForms/internal/interfaces/http/system"

	"github.com/go-chi/chi/v5"
)

func registerRoutes(r *chi.Mux, deps *HTTPDeps) {
	registerSystemRoutes(r, deps.SystemHandler, deps.AuthMiddleware)
	registerProjectRoutes(r, deps.ProjectsHandler, deps.AuthMiddleware)
	registerApiKeyRoutes(r, deps.ApiKeysHandler, deps.AuthMiddleware)
	registerFormRoutes(r, deps.FormsHandler, deps.AuthMiddleware)
}

func registerSystemRoutes(
	r *chi.Mux,
	h *system.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Get("/health", h.Health)
		r.With(authMW.Auth()).
			Get("/protected/health", h.ProtectedHealth)
	})
}

func registerProjectRoutes(
	r *chi.Mux,
	h *projects.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Get("/projects", h.List)
		r.Post("/projects", h.Create)
	})
}

func registerApiKeyRoutes(
	r *chi.Mux,
	h *keys.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Get("/projects/{project_id}/keys", h.List)
		r.Post("/projects/{project_id}/keys", h.Create)
		r.Delete("/projects/{project_id}/keys", h.Revoke)
	})
}

func registerFormRoutes(
	r *chi.Mux,
	h *forms.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Get("/projects/{project_id}/forms", h.List)
		r.Post("/projects/{project_id}/forms", h.Create)
	})
}
