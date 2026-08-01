package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) GetByEmail(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	projectID, err := req.Path("project_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	email, err := req.Path("actor_email").StringRequired()
	if fun.Bail(w, err) {
		return
	}
	members, err := h.ops.GetByEmail(r.Context(), email, projectID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, members)
}
