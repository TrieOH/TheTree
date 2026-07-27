package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (handler *Handlers) CreateOccurrence(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	programID, err := req.Path("program_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.CreateProgramOccurrenceRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	occurrence, err := handler.commands.CreateOccurrence(r.Context(), payload.ToInput(programID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, occurrence, http.StatusCreated)
}
