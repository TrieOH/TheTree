package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) Refresh(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	refreshToken := req.Header("Refresh-Token").String()
	tokens, err := h.commands.Refresh(r.Context(), refreshToken)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, tokens)
}
