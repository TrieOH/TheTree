package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (handler *Handlers) Patch(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	id, err := req.Path("ticket_type_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.PatchTicketTypeRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	ticketType, err := handler.commands.Patch(r.Context(), payload.ToInput(id))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, ticketType)
}
