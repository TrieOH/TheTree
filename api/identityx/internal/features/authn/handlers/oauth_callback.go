package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	provider := chi.URLParam(r, "provider")
	code := r.URL.Query().Get("code")
	if code == "" {
		fun.BadRequest("missing code").Send(w)
		return
	}
	tokens, err := h.commands.OAuthCallback(r.Context(), provider, code)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, tokens, http.StatusCreated)
}
