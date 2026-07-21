package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	"github.com/google/uuid"
)

func (h *Handlers) UpsertSchema(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)

	var projectID *uuid.UUID
	pid, err := req.Path("project_id").UUID()
	if err == nil {
		projectID = &pid
	}

	var payload models.UpsertProfileSchemaRequest
	if bind.BailInto(w, req, &payload) {
		return
	}

	schema, err := h.commands.UpsertSchema(r.Context(), payload.ToInput(projectID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, schema)
}
