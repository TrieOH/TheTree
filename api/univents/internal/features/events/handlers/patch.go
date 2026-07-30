package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (handler *Handlers) Patch(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	eventID, err := req.Path("event_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.PatchEventRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	event, err := handler.commands.Patch(r.Context(), payload.ToInput(eventID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, event)
}
