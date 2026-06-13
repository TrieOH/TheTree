package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/go-chi/chi/v5"
)

// OAuthCallback godoc
// @Summary handles provider OAuth callback
// @Tags authn
// @ID authn_oauth_callback
// @Param provider path string true "Provider" Enums(google, github)
// @Param code query string true "Authorization code"
// @Success 200 {object} fun.Response{data=models.UserTokensOutput}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /auth/{provider}/callback [get]
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
