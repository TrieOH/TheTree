package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) Callback(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	provider, err := req.Path("provider").StringRequired()
	if fun.Bail(w, err) {
		return
	}
	code, err := req.Query("code").StringRequired()
	if fun.Bail(w, err) {
		return
	}
	state, err := req.Query("state").StringRequired()
	if fun.Bail(w, err) {
		return
	}
	finalRedirectURI, err := h.commands.Callback(r.Context(), provider, code, state)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, finalRedirectURI)
}
