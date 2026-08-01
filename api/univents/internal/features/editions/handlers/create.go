package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	eventID, err := req.Path("event_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.CreateEditionRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	edition, err := h.ops.Create(r.Context(), payload.ToInput(eventID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, edition, http.StatusCreated)
}
