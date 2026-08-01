package handlers

import (
	"net/http"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (h *Handlers) PatchProgram(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	id, err := req.Path("program_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.PatchProgramRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	program, err := h.ops.PatchProgram(r.Context(), payload.ToInput(id))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, program)
}
