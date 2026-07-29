package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (handler *Handlers) RevokeSignature(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	token, err := req.Query("token").StringRequired()
	if fun.Bail(w, err) {
		return
	}
	err = handler.commands.RevokeSignature(r.Context(), token)
	if fun.Bail(w, err) {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
