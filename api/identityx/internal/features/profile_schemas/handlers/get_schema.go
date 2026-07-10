package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (h *Handlers) GetSchema(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)

	var projectID *uuid.UUID
	if pid, err := req.Path("project_id").UUID(); err == nil {
		projectID = &pid
	}

	schema, err := h.queries.GetSchema(r.Context(), projectID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, schema)
}
