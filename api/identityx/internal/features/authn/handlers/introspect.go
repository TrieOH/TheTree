package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

// Introspect godoc
// @Summary Introspect current identity
// @Description Returns the identity associated with the current access token.
// @Tags authn
// @ID authn_introspect
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Identity "Current identity"
// @Failure 401 {object} fun.Response "Unauthorized"
// @Failure 500 {object} fun.Response "Internal Server Error"
// @Failure 503 {object} fun.Response "Internal Server Error"
// @Router /auth/introspect [get]
func (h *Handlers) Introspect(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	identity, err := models.RequireIdentity(r.Context())
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, identity)
}
