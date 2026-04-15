package router

import (
	"IdentityX/internal/features/api_keys"
	"IdentityX/internal/features/projects"
	"IdentityX/internal/features/sessions"
	"IdentityX/internal/features/users"
	"IdentityX/internal/interfaces/http/middleware"
	"IdentityX/internal/interfaces/http/system"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/spf13/viper"
)

func registerRoutes(handlers Handlers, r *chi.Mux) *chi.Mux {
	registerAuthRoutes(r, handlers.Users, &handlers.AuthMW)
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
	h *users.Handler,
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

		r.Post("/auth/exchange", h.Exchange)
		r.Post("/auth/refresh", h.Refresh)
		r.With(authMW.Auth(), middleware.NoApiKeys()).
			Post("/auth/logout", h.Logout)
		r.With(httprate.Limit(5, 1*time.Minute, httprate.WithKeyFuncs(httprate.KeyByRealIP))).
			Post("/auth/forgot-password", h.ForgotPassword)
		r.With(httprate.Limit(5, 1*time.Minute, httprate.WithKeyFuncs(httprate.KeyByRealIP))).
			With(middleware.RequireQueryParams("token")).
			Post("/auth/reset-password", h.ResetPassword)
		r.With(authMW.Auth(), middleware.NoApiKeys()).
			With(middleware.RequireQueryParams("token")).
			Post("/auth/verify", h.Verify)
		r.With(authMW.Auth(), middleware.NoApiKeys()).
			Post("/auth/verify/resend", h.ResendVerificationEmail)

		r.Get("/.well-known/jwks.json", h.GetJWKS)

		// FIXME: Create another endpoint for the register that contains SchemaID
		r.With(
			middleware.DefaultQueryParam("schema_type", "core"),
			middleware.DefaultQueryParam("flow_id", "none"),
			middleware.DefaultQueryParam("version", "0"),
		).Post("/projects/{project_id}/register", h.ProjectRegister)

		/*r.With(middleware.DefaultQueryParam("version", "0")).
		Post("/projects/{project_id}/register/{schema_id}", h.ProjectRegister)*/

		r.Post("/projects/{project_id}/logout", h.ProjectLogout)

		if !viper.GetBool("DISABLE_RATE_LIMIT") {
			r.With(httprate.Limit(5, 1*time.Minute, httprate.WithKeyFuncs(httprate.KeyByRealIP))).
				Post("/projects/{project_id}/login", h.ProjectLogin)
		} else {
			r.Post("/projects/{project_id}/login", h.ProjectLogin)
		}
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
		r.Get("/sessions", h.ListUserSessions)
		r.Get("/sessions/me", h.Me)
		r.Delete("/sessions/{session_id}", h.RevokeUserSessionByID)
		r.Delete("/sessions/others", h.RevokeOtherSessions)
		r.Delete("/sessions", h.RevokeAllSessions)
	})
}

func registerProjectRoutes(
	r *chi.Mux,
	h *projects.Handler,
	authMW *middleware.AuthMiddleware,
) {
	r.With(authMW.Auth(), middleware.ClientOnly()).Group(func(r chi.Router) {
		r.Get("/projects/{project_id}/.well-known/jwks.json", h.GetProjectJWKS)
		r.Post("/projects", h.CreateProject)
		r.Get("/projects", h.ListProjects)
		r.Get("/projects/{project_id}", h.GetProjectByID)
		r.Patch("/projects/{project_id}", h.UpdateProjectByID)
		r.Delete("/projects/{project_id}", h.DeleteProjectByID)
		r.Get("/projects/{project_id}/users", h.ListProjectUsers)
		r.Get("/projects/{project_id}/users/{user_id}", h.GetProjectUserByID)
	})
}
