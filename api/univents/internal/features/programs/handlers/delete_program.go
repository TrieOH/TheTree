package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) DeleteProgram(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	id, err := req.Path("program_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	program, err := h.ops.DeleteProgram(r.Context(), id)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, program)
}
