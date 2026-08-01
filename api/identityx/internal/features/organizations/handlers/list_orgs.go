package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) ListOrgs(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	namespaces, err := h.ops.ListOrgs(r.Context())
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, namespaces)
}
