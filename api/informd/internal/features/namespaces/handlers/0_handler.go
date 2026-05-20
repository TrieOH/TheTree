package handlers

import (
	"Informd/internal/features/namespaces/commands"
	"Informd/internal/features/namespaces/queries"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	commands *commands.CommandService
	queries  *queries.QueryService
}

func NewHandler(
	commands *commands.CommandService,
	queries *queries.QueryService,
) *Handlers {
	return &Handlers{
		commands: commands,
		queries:  queries,
	}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	jwt func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(jwt)
		r.Get("/namespaces", h.ListNamespaces)
		r.Post("/namespaces", h.Create)
		r.Get("/namespaces/{namespace_id}/members", h.ListMembers)
		r.Post("/namespaces/{namespace_id}/members", h.AddMember)
		r.Delete("/namespaces/{namespace_id}/members", h.RemoveMember)
		r.Post("/namespaces/{namespace_id}/forms", h.CreateForm)
		r.Get("/namespaces/{namespace_id}/forms", h.ListForms)
		r.Get("/namespaces/{namespace_id}/forms/{form_id}/members", h.ListFormMembers)
		r.Post("/namespaces/{namespace_id}/forms/{form_id}/members", h.AddFormMember)
		r.Delete("/namespaces/{namespace_id}/forms/{form_id}/members", h.RemoveFormMember)
	})
}
