package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (handler *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	var payload models.CreateEventRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	event, err := handler.commands.Create(r.Context(), payload)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, event, http.StatusCreated)
}
