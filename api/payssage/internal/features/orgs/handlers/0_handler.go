package handlers

import (
	"net/http"
	"payssage/internal/features/orgs"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *orgs.Operations
}

func NewHandlers(ops *orgs.Operations) *Handlers {
	return &Handlers{ops: ops}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	jwtAuth func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(jwtAuth)
		r.Get("/organizations", h.ListOrgs)
		r.Post("/organizations", h.Create)
		r.Get("/organizations/{organization_id}/members", h.ListMembers)
		r.Post("/organizations/{organization_id}/members", h.AddMember)
		r.Delete("/organizations/{organization_id}/members", h.RemoveMember)
		r.Get("/organizations/{organization_id}/member/{member_id}", h.GetMemberByID)
		r.Get("/organizations/{organization_id}/member/{member_email}:by_email", h.GetMemberByEmail)
	})
}
