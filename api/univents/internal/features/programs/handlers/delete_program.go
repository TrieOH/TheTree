package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) DeleteProgram(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	id, err := req.Path("program_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	program, err := handler.commands.DeleteProgram(r.Context(), id)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, program)
}
