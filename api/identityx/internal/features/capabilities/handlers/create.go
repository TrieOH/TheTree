package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// Create godoc
// @Summary Create a capability in a project
// @Tags capabilities
// @ID capabilities_create
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param project_id path uuid.UUID true "Project ID"
// @Param request body models.CreateCapabilityRequest true "Capability creation data"
// @Success 200 {object} fun.Response{data=models.Capability} "Capability data"
// @Failure 401 {object} fun.Response "Unauthorized"
// @Failure 404 {object} fun.Response "Bad Request"
// @Failure 500 {object} fun.Response "Internal Server Error"
// @Failure 503 {object} fun.Response "Internal Server Error"
// @Router /projects/{project_id}/capabilities [post]
func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	projectID, err := req.Path("project_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.CreateCapabilityRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	capability, err := h.commands.Create(r.Context(), payload.ToInput(projectID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, capability, http.StatusCreated)
}
