package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// AddMember godoc
// @Summary Add a project member
// @Description Lets you add a member to the project as a member or admin
// @Tags projects
// @ID projects_addmember
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.AddProjectMemberRequest true "Member details"
// @Success 201 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /projects/{project_id}/members [post]
func (h *Handlers) AddMember(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	projectID, err := req.Path("project_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.AddProjectMemberRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err = h.commands.AddMember(r.Context(), payload.ToInput(projectID))
	if fun.Bail(w, err) {
		return
	}
	fun.Created().Send(w)
}
