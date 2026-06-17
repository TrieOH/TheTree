package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

// IsSetup godoc
// @Summary gets setup status
// @Tags authn
// @ID authn_issetup
// @Produce json
// @Success 200 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /auth/setup [get]
func (h *Handlers) IsSetup(w http.ResponseWriter, _ *http.Request) {
	if globals.SetupComplete() {
		fun.ServiceUnavailable("setup already complete").Send(w)
		return
	}
	fun.OK().Send(w)
}
