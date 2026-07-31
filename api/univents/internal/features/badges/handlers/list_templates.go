package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// ListTemplates lists all badge templates for an edition
//
//	@Summary		List badge templates
//	@Description	List all badge templates for an edition
//	@Tags			badges
//	@Produce		json
//	@Param			event_id	path		string								true	"Event ID"
//	@Param			edition_id	path		string								true	"Edition ID"
//	@Success		200			{array}		models.BadgeTemplate
//	@Failure		400			{object}	errx.ErrorResponse
//	@Failure		401			{object}	errx.ErrorResponse
//	@Failure		500			{object}	errx.ErrorResponse
//	@Security		BearerAuth
//	@Router			/events/{event_id}/editions/{edition_id}/badges [get]
func (h *Handler) ListTemplates(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}

	templates, err := h.quers.ListTemplates(r.Context(), editionID)
	if fun.Bail(w, err) {
		return
	}

	fun.Respond(w, templates, http.StatusOK)
}
