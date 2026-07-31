package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) GetByID(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	id, err := req.Path("signature_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	signature, err := handler.queries.GetByID(r.Context(), id)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, signature)
}
