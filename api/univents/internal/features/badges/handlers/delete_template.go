package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// DeleteTemplate deletes a badge template
//
//	@Summary		Delete badge template
//	@Description	Delete a badge template
//	@Tags			badges
//	@Param			event_id	path		string								true	"Event ID"
//	@Param			edition_id	path		string								true	"Edition ID"
//	@Param			template_id	path		string								true	"Template ID"
//	@Success		204			"No Content"
//	@Failure		400			{object}	errx.ErrorResponse
//	@Failure		401			{object}	errx.ErrorResponse
//	@Failure		404			{object}	errx.ErrorResponse
//	@Failure		500			{object}	errx.ErrorResponse
//	@Security		BearerAuth
//	@Router			/events/{event_id}/editions/{edition_id}/badges/{template_id} [delete]
func (h *Handler) DeleteTemplate(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	templateID, err := req.Path("template_id").UUID()
	if fun.Bail(w, err) {
		return
	}

	err = h.cmds.DeleteTemplate(r.Context(), editionID, templateID)
	if fun.Bail(w, err) {
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
