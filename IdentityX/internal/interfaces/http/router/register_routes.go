package router

import (
	"IdentityX/internal/features/account"
	"IdentityX/internal/features/api_keys"
	"IdentityX/internal/features/auth"
	"IdentityX/internal/features/projects"
	"IdentityX/internal/features/sessions"
	"IdentityX/internal/interfaces/http/middleware"
	"IdentityX/internal/interfaces/http/system"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/spf13/viper"
)

func registerRoutes(handlers Handlers, r *chi.Mux) *chi.Mux {
	registerAuthRoutes(r, handlers.Users, &handlers.AuthMW)
	registerAccountRoutes(r, handlers.Accounts, &handlers.AuthMW)
	registerSessionRoutes(r, handlers.Sessions, &handlers.AuthMW)
	registerProjectRoutes(r, handlers.Projects, &handlers.AuthMW)
	registerApiKeyRoutes(r, handlers.ApiKeys, &handlers.AuthMW)
	registerSystemRoutes(r, handlers.System, &handlers.AuthMW)

	return r
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

func registerApiKeyRoutes(
	r *chi.Mux,
	h *api_keys.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Use(middleware.ClientOnly())
		r.Post("/projects/{project_id}/api-keys/rotate", h.RotateApiKey)
		r.Delete("/projects/{project_id}/api-keys", h.RevokeApiKey)
	})
}

func registerAuthRoutes(
	r *chi.Mux,
	h *auth.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Post("/auth/register", h.Register)

		if !viper.GetBool("DISABLE_RATE_LIMIT") {
			r.With(httprate.Limit(5, 1*time.Minute, httprate.WithKeyFuncs(httprate.KeyByRealIP))).
				Post("/auth/login", h.Login)
		} else {
			r.Post("/auth/login", h.Login)
		}

		r.Post("/auth/refresh", h.Refresh)
		r.With(authMW.Auth(), middleware.NoApiKeys()).
			Post("/auth/logout", h.Logout)

		r.With(middleware.AllowOnlyQueryParams("project_id")).
			Get("/.well-known/jwks.json", h.GetJWKS)

		r.Post("/projects/{project_id}/register", h.ProjectRegister)

		if !viper.GetBool("DISABLE_RATE_LIMIT") {
			r.With(httprate.Limit(5, 1*time.Minute, httprate.WithKeyFuncs(httprate.KeyByRealIP))).
				Post("/projects/{project_id}/login", h.ProjectLogin)
		} else {
			r.Post("/projects/{project_id}/login", h.ProjectLogin)
		}
	})
}

func registerAccountRoutes(
	r *chi.Mux,
	h *account.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(httprate.Limit(5, 1*time.Minute, httprate.WithKeyFuncs(httprate.KeyByRealIP)))
		r.Use(middleware.NoApiKeys())

		r.Post("/account/forgot-password", h.ForgotPassword)
		r.With(authMW.Auth()).Post("/account/verify/resend", h.ResendVerificationEmail)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireQueryParams("token"))

			r.Post("/account/reset-password", h.ResetPassword)
			r.With(authMW.Auth()).Post("/account/verify", h.Verify)
		})
	})
}

func registerSessionRoutes(
	r *chi.Mux,
	h *sessions.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Use(middleware.NoApiKeys())
		r.Get("/sessions", h.List)
		r.Get("/sessions/me", h.Me)
		r.Delete("/sessions/{session_id}", h.RevokeByID)
		r.Delete("/sessions/others", h.RevokeOthers)
		r.Delete("/sessions", h.RevokeAll)
	})
}

func registerProjectRoutes(
	r *chi.Mux,
	h *projects.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.With(authMW.Auth(), middleware.ClientOnly()).Group(func(r chi.Router) {
		r.Post("/projects", h.Create)
		r.Get("/projects", h.List)
		r.Get("/projects/{project_id}", h.GetByID)
		r.Patch("/projects/{project_id}", h.Update)
		r.Delete("/projects/{project_id}", h.Delete)
		r.Get("/projects/{project_id}/users", h.ListProjectUsers)
		r.Get("/projects/{project_id}/users/{user_id}", h.GetProjectUserByID)
	})
}
