package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// CreateTemplate creates a new badge template
//
//	@Summary		Create badge template
//	@Description	Create a new badge template for an edition
//	@Tags			badges
//	@Accept			json
//	@Produce		json
//	@Param			event_id	path		string								true	"Event ID"
//	@Param			edition_id	path		string								true	"Edition ID"
//	@Param			body		body		models.CreateBadgeTemplateRequest	true	"Badge Template"
//	@Success		201			{object}	models.BadgeTemplate
//	@Failure		400			{object}	errx.ErrorResponse
//	@Failure		401			{object}	errx.ErrorResponse
//	@Failure		500			{object}	errx.ErrorResponse
//	@Security		BearerAuth
//	@Router			/events/{event_id}/editions/{edition_id}/badges [post]
func (h *Handler) CreateTemplate(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	
	var payload models.CreateBadgeTemplateRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	
	input := models.CreateBadgeTemplateInput{
		EditionID:    editionID,
		TicketTypeID: payload.TicketTypeID,
		Name:         payload.Name,
		DesignData:   payload.DesignData,
	}
	
	template, err := h.cmds.CreateTemplate(r.Context(), input)
	if fun.Bail(w, err) {
		return
	}
	
	fun.Respond(w, template, http.StatusCreated)
}
