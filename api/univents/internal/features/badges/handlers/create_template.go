package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

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

	input := payload.ToInput(editionID)

	template, err := h.commands.CreateTemplate(r.Context(), input)
	if fun.Bail(w, err) {
		return
	}

	fun.Respond(w, template, http.StatusCreated)
}
