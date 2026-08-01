package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	var payload models.CreateEventRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	event, err := h.ops.Create(r.Context(), payload)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, event, http.StatusCreated)
}
