package router

import (
	"GoAuth/internal/adapters/http/handlers"
	"GoAuth/internal/adapters/http/middleware"
	"GoAuth/internal/adapters/observability/logs"
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/application"
	"GoAuth/internal/infrastructure"
	"GoAuth/internal/infrastructure/telemetry"
	"database/sql"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
)

func registerRoutes(db *sql.DB, r *chi.Mux) *chi.Mux {
	queries := sqlc.New(db)
	txRunner := transactions.NewTxRunner(db)
	tracer := otel.Tracer(string(telemetry.GoAuthTracer))
	infra := infrastructure.NewInfra(db, queries, txRunner, logs.L(), tracer)

	app := application.NewApplication(infra)

	handlerBundle := handlers.New(app)

	authMW := middleware.NewAuthMiddleware(app.Authenticator, tracer, viper.GetString("ISSUER"))

	registerAuthRoutes(r, handlerBundle.AuthHandler, authMW)
	registerSessionRoutes(r, handlerBundle.SessionHandler, authMW)
	registerProjectRoutes(r, handlerBundle.ProjectHandler, authMW)
	registerSchemaRoutes(r, handlerBundle.SchemaHandler, authMW)
	registerSchemaVersionRoutes(r, handlerBundle.SchemaVersionHandler, authMW)
	registerSchemaFieldsRoutes(r, handlerBundle.SchemaFieldsHandler, authMW)
	registerScopeRoutes(r, handlerBundle.ScopeHandler, authMW)
	registerPermissionRoutes(r, handlerBundle.PermissionHandler, authMW)
	registerRoleRoutes(r, handlerBundle.RoleHandler, authMW)

	return r
}

func registerAuthRoutes(
	r *chi.Mux,
	h *handlers.AuthHandler,
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
		r.With(authMW.Auth()).
			Post("/auth/logout", h.Logout)
		r.With(authMW.Auth()).
			With(middleware.RequireQueryParams("token")).
			Post("/auth/verify", h.Verify)
		r.With(authMW.Auth()).
			Post("/auth/verify/resend", h.ResendVerificationEmail)

		r.Get("/.well-known/jwks.json", h.GetJWKS)

		// FIXME: Create another endpoint for the register that contains SchemaID
		r.With(
			middleware.DefaultQueryParam("schema_type", "core"),
			middleware.DefaultQueryParam("flow_id", "none"),
			middleware.DefaultQueryParam("version", "0"),
		).Post("/projects/{project_id}/register", h.ProjectRegister)

		/*r.With(
			middleware.DefaultQueryParam("schema_type", "core"),
			middleware.DefaultQueryParam("flow_id", "none"),
			middleware.DefaultQueryParam("version", "0"),
		).Post("/projects/{project_id}/register/{schema_id}", h.ProjectRegister)*/

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
	h *handlers.SessionHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Get("/sessions", h.ListUserSessions)
		r.Get("/sessions/me", h.Me)
		r.Delete("/sessions/{session_id}", h.RevokeUserSessionByID)
		r.Delete("/sessions/others", h.RevokeOtherSessions)
		r.Delete("/sessions", h.RevokeAllSessions)
	})
}

func registerProjectRoutes(
	r *chi.Mux,
	h *handlers.ProjectHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Get("/projects/{project_id}/.well-known/jwks.json", h.GetProjectJWKS)
		r.With(authMW.Auth(), middleware.ClientOnly()).Group(func(r chi.Router) {
			r.Post("/projects", h.CreateProject)
			r.Get("/projects", h.ListProjects)
			r.Get("/projects/{project_id}", h.GetProjectByID)
			r.Patch("/projects/{project_id}", h.UpdateProjectByID)
			r.Delete("/projects/{project_id}", h.DeleteProjectByID)
		})
	})
}

func registerSchemaRoutes(
	r *chi.Mux,
	schemas *handlers.SchemaHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Use(middleware.ClientOnly())
		r.Route("/projects/{project_id}/schemas", func(r chi.Router) {
			r.Post("/", schemas.Draft)
			r.Get("/", schemas.List)
			r.Get("/ids", schemas.GetIDsFromProjectID)
			r.Get("/{schema_id}", schemas.GetByID)
			/* r.With(
				middleware.DefaultQueryParam("schema_type", "context"),
				middleware.DefaultQueryParam("flow_id", "none"),
			).Get("/{schema_id}/latest", schemas.GetLatestForm)
			r.With(
				middleware.DefaultQueryParam("schema_type", "context"),
				middleware.DefaultQueryParam("flow_id", "none"),
			).Get("/{schema_id}/v{version:[0-9]+}", schemas.GetSpecificForm) */
			r.Get("/{schema_id}/verbose", schemas.GetVerbose)
			r.Post("/{schema_id}/publish", schemas.Publish)
		})
	})
}

func registerSchemaVersionRoutes(
	r *chi.Mux,
	h *handlers.SchemaVersionHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Use(middleware.ClientOnly())
		r.Route("/projects/{project_id}/schemas/{schema_id}/versions", func(r chi.Router) {
			r.Post("/draft", h.Draft)
			r.Post("/publish", h.Publish)
			r.Get("/current", h.GetCurrent)
			r.Get("/latest", h.GetLatest)
			r.Get("/v{version:[0-9]+}", h.GetVerbose)
		})
	})
}

func registerSchemaFieldsRoutes(
	r *chi.Mux,
	h *handlers.SchemaFieldsHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Use(middleware.ClientOnly())
		r.Post("/projects/{project_id}/schemas/{schema_id}/v{version:[0-9]+}", h.Create)
	})
}

func registerScopeRoutes(
	r *chi.Mux,
	h *handlers.ScopeHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Use(middleware.ClientOnly())
		r.Route("/projects/{project_id}/scopes", func(r chi.Router) {
			r.Post("/", h.Create)
			r.Get("/", h.GetProjectScopes)
			r.Get("/{scope_id}", h.GetByID)
		})
	})
}

func registerPermissionRoutes(
	r *chi.Mux,
	h *handlers.PermissionHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Use(middleware.ClientOnly())
		r.With(middleware.AllowOnlyQueryParams("scope_id")).
			Get("/projects/{project_id}/identities/{entity_id}/permissions", h.GetEffective)
		r.Post("/projects/{project_id}/identities/{entity_id}/permissions", h.GiveDirect)
		r.Delete("/projects/{project_id}/identities/{entity_id}/permissions", h.TakeDirect)
		r.Route("/projects/{project_id}/permissions", func(r chi.Router) {
			r.Post("/", h.Create)
			r.Get("/{permission_id}", h.GetByID)
			r.With(middleware.AllowOnlyQueryParams("object", "action")).
				Get("/", h.ListByProject)
		})
	})
}

func registerRoleRoutes(
	r *chi.Mux,
	h *handlers.RoleHandler,
	authMW *middleware.AuthMiddleware,
) {
	r.Group(func(r chi.Router) {
		r.Use(authMW.Auth())
		r.Use(middleware.ClientOnly())

		r.Post("/projects/{project_id}/roles", h.Create)
		r.Get("/projects/{project_id}/roles/{role_id}", h.GetByID)
		r.Patch("/projects/{project_id}/roles/{role_id}", h.UpdateDescription)
		r.Get("/projects/{project_id}/roles", h.ListByProject)
		r.With(middleware.RequireOnlyQueryParams("name")).
			Get("/projects/{project_id}/roles/search", h.GetByName)

		r.Get("/projects/{project_id}/roles/{role_id}/permissions", h.GetPermissions)
		r.Post("/projects/{project_id}/roles/{role_id}/permissions/{permission_id}", h.AddPermission)
		r.Delete("/projects/{project_id}/roles/{role_id}/permissions/{permission_id}", h.RemovePermission)

		r.Get("/projects/{project_id}/identities/{entity_id}/roles", h.GetUserRoles)
		r.Post("/projects/{project_id}/identities/{entity_id}/roles", h.GiveRole)
		r.Delete("/projects/{project_id}/identities/{entity_id}/roles", h.TakeRole)
	})
}
