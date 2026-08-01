package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (h *Handlers) AddMember(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	eventID, err := req.Path("event_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.AddEventMemberRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	member, err := h.ops.AddMember(r.Context(), eventID, payload)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, member, http.StatusCreated)
}
