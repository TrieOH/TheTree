package handlers

import (
	"IdentityX/internal/features/authn/commands"
	"IdentityX/internal/features/authn/queries"
	"IdentityX/models"
	"net/http"

	"github.com/MintzyG/fun/middlewares"
	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	commands *commands.Commands
	queries  *queries.Queries
}

func NewHandlers(
	commands *commands.Commands,
	queries *queries.Queries,
) *Handlers {
	return &Handlers{
		commands: commands,
		queries:  queries,
	}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	jwtAuth func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Post("/auth/setup", h.Setup)
		r.Get("/auth/setup", h.IsSetup)
		r.With(middlewares.WithParams[models.ProjectIDQueryParam]()).Post("/auth/register", h.Register)
		r.With(middlewares.WithParams[models.ProjectIDQueryParam]()).Post("/auth/login", h.Login)
		r.With(jwtAuth).Post("/auth/logout", h.Logout)
		r.Post("/auth/refresh", h.Refresh)
		r.With(middlewares.WithParams[models.ProjectIDQueryParam]()).Get("/auth/{provider}/connect", h.OAuthConnect)
		r.Get("/auth/{provider}/callback", h.OAuthCallback)
		r.With(middlewares.WithParams[models.ProjectIDQueryParam]()).Get("/.well-known/jwks.json", h.JWKS)
		r.With(jwtAuth).Get("/auth/introspect", h.Introspect)
	})
}
