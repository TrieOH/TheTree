package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (handler *Handlers) CreateProgram(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	editionID, err := req.Path("edition_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.CreateProgramRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	program, err := handler.commands.CreateProgram(r.Context(), payload.ToInput(editionID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, program, http.StatusCreated)
}
