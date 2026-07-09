package handlers

import (
	"net/http"
	"univents/contracts"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (handler *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload contracts.CreateActivityRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	activity, err := handler.commands.Create(r.Context(), payload.ToSpec(editionID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, activity, http.StatusCreated)
}
