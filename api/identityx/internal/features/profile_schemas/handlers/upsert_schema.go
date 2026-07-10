package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	"github.com/google/uuid"
)

// UpsertSchema godoc
// @Summary Upsert the profile schema for a project
// @Tags profile_schemas
// @ID profile_schemas_upsertprojectschema
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Param body body models.UpsertProfileSchemaRequest true "Schema payload"
// @Success 200 {object} fun.Response{data=models.ProjectProfileSchema}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /projects/{project_id}/profile-schema [put]

// UpsertPlatformSchema godoc
// @Summary Upsert the platform-wide profile schema
// @Tags profile_schemas
// @ID profile_schemas_upsertplatformschema
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.UpsertProfileSchemaRequest true "Schema payload"
// @Success 200 {object} fun.Response{data=models.ProjectProfileSchema}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /profile-schema [put]
func (h *Handlers) UpsertSchema(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)

	var projectID *uuid.UUID
	if pid, err := req.Path("project_id").UUID(); err == nil {
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
