package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.CreateTicketTypeRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	ticketType, err := h.ops.Create(r.Context(), payload.ToInput(editionID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, ticketType, http.StatusCreated)
}
