package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// GetSchema godoc
// @Summary Get the profile schema for a project
// @Tags profile_schemas
// @ID profile_schemas_getprojectschema
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path string true "Project ID"
// @Success 200 {object} fun.Response{data=models.ProjectProfileSchema}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /projects/{project_id}/profile-schema [get]

// GetPlatformSchema godoc
// @Summary Get the platform-wide profile schema
// @Tags profile_schemas
// @ID profile_schemas_getplatformschema
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fun.Response{data=models.ProjectProfileSchema}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /profile-schema [get]
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
