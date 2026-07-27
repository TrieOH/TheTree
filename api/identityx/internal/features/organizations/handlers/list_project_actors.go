package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

func (h *Handlers) ListProjectActors(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	projectID, err := req.Path("project_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	orgID, err := req.Path("org_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	actors, err := h.queries.ListProjectActors(r.Context(), orgID, projectID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, actors)
}
