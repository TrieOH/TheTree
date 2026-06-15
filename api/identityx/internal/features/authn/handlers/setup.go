package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// Setup godoc
// @Summary Setups IDX
// @Description Creates the first account in the system as a super admin and enables authn
// @Tags authn
// @ID authn_setup
// @Accept json
// @Produce json
// @Param request body models.IDXLoginRequest true "setup details"
// @Success 201 {object} fun.Response{data=models.UserTokensOutput}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /auth/setup [post]
func (h *Handlers) Setup(w http.ResponseWriter, r *http.Request) {
	if globals.SetupComplete() {
		fun.Forbidden("setup already complete").Send(w)
		return
	}
	req := fun.From(r)
	ctx := r.Context()
	var payload models.IDXLoginRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err := h.commands.Setup(ctx, payload.ToSetupInput())
	if fun.Bail(w, err) {
		return
	}
	tokens, err := h.commands.Login(ctx, payload.ToInput(nil))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, tokens, http.StatusCreated)
	globals.MarkSetupComplete()
}
