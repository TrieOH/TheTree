package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// Login godoc
// @Summary logins a user on IDX
// @Description This route is disabled until setup is complete
// @Tags authn
// @ID authn_login
// @Accept json
// @Produce json
// @Param request body models.IDXLoginRequest true "login details"
// @Success 200 {object} fun.Response{data=models.UserTokensOutput}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /auth/login [post]
func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	var payload models.IDXLoginRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	tokens, err := h.commands.Login(r.Context(), payload.ToInput())
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, tokens)
}
