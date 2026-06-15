package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/go-chi/chi/v5"
)

// OAuthConnect godoc
// @Summary gets URL for provider consent screen
// @Tags authn
// @ID authn_oauth_connect
// @Param provider path string true "Provider" Enums(google, github)
// @Success 200 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /auth/{provider}/connect [get]
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
