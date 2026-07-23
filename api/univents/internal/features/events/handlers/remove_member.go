package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (handler *Handlers) RemoveMember(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	eventID, err := req.Path("event_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.RemoveMemberRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err = handler.commands.RemoveMember(r.Context(), eventID, payload)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, nil, http.StatusNoContent)
}
