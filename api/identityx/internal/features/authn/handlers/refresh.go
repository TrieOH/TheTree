package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

// Refresh godoc
// @Summary refreshes IDX tokens
// @Description takes a valid refresh token, invalidates both access and refresh and issue new tokens.
// @Tags authn
// @ID authn_refresh
// @Accept json
// @Produce json
// @Success 200 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /auth/refresh [post]
func (h *Handlers) Refresh(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	refreshToken := req.Header("refresh_token").String()
	tokens, err := h.commands.Refresh(r.Context(), refreshToken)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, tokens)
}
