package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// GetTemplate gets a badge template by ID
//
//	@Summary		Get badge template
//	@Description	Get a badge template by ID
//	@Tags			badges
//	@Produce		json
//	@Param			event_id	path		string								true	"Event ID"
//	@Param			edition_id	path		string								true	"Edition ID"
//	@Param			template_id	path		string								true	"Template ID"
//	@Success		200			{object}	models.BadgeTemplate
//	@Failure		400			{object}	errx.ErrorResponse
//	@Failure		401			{object}	errx.ErrorResponse
//	@Failure		404			{object}	errx.ErrorResponse
//	@Failure		500			{object}	errx.ErrorResponse
//	@Security		BearerAuth
//	@Router			/events/{event_id}/editions/{edition_id}/badges/{template_id} [get]
func (h *Handler) GetTemplate(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	templateID, err := req.Path("template_id").UUID()
	if fun.Bail(w, err) {
		return
	}

	template, err := h.quers.GetTemplate(r.Context(), editionID, templateID)
	if fun.Bail(w, err) {
		return
	}

	fun.Respond(w, template, http.StatusOK)
}
