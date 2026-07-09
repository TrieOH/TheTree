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
	jwtAuth func(http.Handler) http.Handler,
) {
	r.Group(func(r chi.Router) {
		r.Get("/verify/{hash}", h.Verify)
		r.With(jwtAuth).Group(func(r chi.Router) {
			// Templates
			r.Post("/events/{event_id}/editions/{edition_id}/certification-templates", h.CreateTemplate)
			r.Get("/events/{event_id}/editions/{edition_id}/certification-templates", h.ListTemplates)
			r.Get("/events/{event_id}/editions/{edition_id}/certification-templates/{template_id}", h.GetTemplateByID)

			// Certifications
			r.Post("/events/{event_id}/editions/{edition_id}/certifications", h.Certify)
			r.Get("/certifications/{cert_id}", h.GetCertByID)
			r.Get("/users/{user_id}/certifications", h.ListByUser)
			r.Get("/certifications", h.ListByTarget)

			// Linking
			r.Patch("/events/{event_id}/editions/{edition_id}/activities/{activity_id}/certification-templates/set", h.SetActivityTemplate)
			r.Patch("/events/{event_id}/editions/{edition_id}/certification-templates/set", h.SetEditionTemplate)
		})

	})
}
