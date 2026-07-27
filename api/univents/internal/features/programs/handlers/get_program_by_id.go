package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) GetProgramByID(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	id, err := req.Path("program_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	program, err := handler.queries.GetProgramByID(r.Context(), id)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, program)
}
