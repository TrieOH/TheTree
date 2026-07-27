package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

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
