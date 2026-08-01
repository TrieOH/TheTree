package handlers

import (
	"IdentityX/internal/features/organizations"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handlers struct {
	ops *organizations.Operations
}

func NewHandlers(ops *organizations.Operations) *Handlers {
	return &Handlers{
		ops: ops,
	}
}

func RegisterRoutes(
	r *chi.Mux,
	h *Handlers,
	jwtAuth func(http.Handler) http.Handler,
	clientOnly func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Use(jwtAuth, clientOnly)
		r.Get("/organizations", h.ListOrgs)
		r.Post("/organizations", h.Create)
		r.Get("/organizations/{organization_id}/members", h.ListMembers)
		r.Post("/organizations/{organization_id}/members", h.AddMember)
		r.Delete("/organizations/{organization_id}/members", h.RemoveMember)
		r.Get("/organizations/{organization_id}/projects", h.ListProjects)
		r.Post("/organizations/{organization_id}/projects", h.CreateProject)
		r.Post("/organizations/{org_id}/projects/{project_id}/actors", h.CreateProjectActor)
		r.Get("/organizations/{org_id}/projects/{project_id}/actors", h.ListProjectActors)
		r.Get("/organizations/{organization_id}/projects/{project_id}/members", h.ListProjectMembers)
		r.Post("/organizations/{organization_id}/projects/{project_id}/members", h.AddProjectMember)
		r.Delete("/organizations/{organization_id}/projects/{project_id}/members", h.RemoveProjectMember)
		r.Get("/organizations/{organization_id}/projects/{project_id}/actors/{actor_id}", h.GetActorByID)
	})
}
