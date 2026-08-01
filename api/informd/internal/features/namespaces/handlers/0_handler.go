package handlers

import (
	"Informd/internal/features/namespaces"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *namespaces.Operations
}

func NewHandler(ops *namespaces.Operations) *Handlers {
	return &Handlers{
		ops: ops,
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
		r.Get("/namespaces/{namespace_id}/forms/archived", h.ListArchivedForms)
		// TODO: kill these duplicated namespaced routes — CheckForm already anchors via the form's namespace.
		r.Get("/namespaces/{namespace_id}/forms/{form_id}/full", h.GetFullFormNamespaced)
		r.Get("/namespaces/{namespace_id}/forms/{form_id}/members", h.ListFormMembers)
		r.Post("/namespaces/{namespace_id}/forms/{form_id}/members", h.AddFormMember)
		r.Delete("/namespaces/{namespace_id}/forms/{form_id}/members", h.RemoveFormMember)
		r.Post("/namespaces/{namespace_id}/forms/{form_id}/open", h.Open)
		r.Post("/namespaces/{namespace_id}/forms/{form_id}/close", h.Close)
		r.Post("/namespaces/{namespace_id}/forms/{form_id}/archive", h.Archive)
		r.Post("/namespaces/{namespace_id}/forms/{form_id}/redraft", h.ReDraft)
		r.Get("/namespaces/{namespace_id}/forms/{form_id}/responses/count", h.ResponseCount)
	})
}
