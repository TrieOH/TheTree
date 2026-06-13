package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// Register godoc
// @Summary registers a user to IDX
// @Description This route is disabled until setup is complete
// @Tags authn
// @ID authn_register
// @Accept json
// @Produce json
// @Param request body models.IDXRegisterRequest true "register details"
// @Success 201 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /auth/register [post]
func (h *Handlers) Register(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	var payload models.IDXRegisterRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err := h.commands.Register(r.Context(), payload.ToInput())
	if fun.Bail(w, err) {
		return
	}
	fun.Created().Send(w)
}
