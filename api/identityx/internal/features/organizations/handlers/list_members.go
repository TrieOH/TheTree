package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) ListMembers(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	orgID, err := req.Path("organization_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	members, err := h.queries.ListMembers(r.Context(), orgID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, members)
}
