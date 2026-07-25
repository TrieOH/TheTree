package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) OAuthConnect(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	provider := chi.URLParam(r, "provider")
	url, err := h.commands.OAuthConnect(r.Context(), provider)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, map[string]string{"url": url})
}
