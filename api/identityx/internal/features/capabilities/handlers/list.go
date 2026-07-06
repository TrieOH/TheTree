package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

// List godoc
// @Summary List capabilities in a project
// @Tags capabilities
// @ID capabilities_list
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path uuid.UUID true "Project ID"
// @Success 200 {array} models.Capability "Capabilties data"
// @Failure 401 {object} fun.Response "Unauthorized"
// @Failure 404 {object} fun.Response "Bad Request"
// @Failure 500 {object} fun.Response "Internal Server Error"
// @Failure 503 {object} fun.Response "Internal Server Error"
// @Router /projects/{project_id}/capabilities [get]
func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	projectID, err := req.Path("project_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	capabilities, err := h.queries.List(r.Context(), projectID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, capabilities)
}
