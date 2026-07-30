package handlers

import (
	"net/http"
	"univents/internal/features/certifications/commands"
	"univents/internal/features/certifications/queries"

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
	jwt func(http.Handler) http.Handler,
) {
	r.Get("/editions/{edition_id}/certifications/templates", h.ListTemplates)
	r.Get("/certifications/templates/{template_id}", h.GetTemplateByID)
	r.With(jwt).Post("/editions/{edition_id}/certifications/templates", h.CreateTemplate)
	r.With(jwt).Put("/certifications/templates/{template_id}", h.UpdateTemplate)
	r.With(jwt).Delete("/certifications/templates/{template_id}", h.DeleteTemplate)
	r.With(jwt).Post("/certifications/templates/{template_id}/link", h.LinkCertTemplate)
	r.With(jwt).Delete("/certifications/templates/{template_id}/link", h.UnlinkCertTemplate)
	r.Get("/certifications/templates/{template_id}/links", h.ListCertTemplateLinks)

	r.Get("/verify/{hash}", h.VerifyCert)

	r.With(jwt).Get("/certifications/{cert_id}", h.GetCertByID)
	r.With(jwt).Get("/editions/{edition_id}/certifications", h.ListCertsByEdition)
	r.With(jwt).Get("/certifications", h.ListMyCerts)
	r.With(jwt).Post("/certifications/{cert_id}/invalidate", h.InvalidateCert)
	r.With(jwt).Get("/editions/{edition_id}/certifications/emission-errors", h.ListEmissionErrors)
}
