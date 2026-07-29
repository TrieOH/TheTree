package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (handler *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.AddSignatureRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	signature, err := handler.commands.Create(r.Context(), payload.ToInput(editionID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, signature, http.StatusCreated)
}
