package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// Create godoc
// @Summary Create a project
// @Description Creates a project with the subject as the owner, only allows IDX Clients
// @Tags projects
// @ID projects_create
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateProjectRequest true "Project creation data"
// @Success 200 {object} fun.Response{data=models.Project} "Project data"
// @Failure 401 {object} fun.Response "Unauthorized"
// @Failure 404 {object} fun.Response "Bad Request"
// @Failure 500 {object} fun.Response "Internal Server Error"
// @Failure 503 {object} fun.Response "Internal Server Error"
// @Router /projects [post]
func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	var payload models.CreateProjectRequest
	if bind.BailInto(w, fun.From(r), &payload) {
		return
	}
	org, err := h.commands.Create(r.Context(), payload.ToInput(nil))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, org, http.StatusCreated)
}
