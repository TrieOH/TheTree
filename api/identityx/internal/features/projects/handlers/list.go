package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

// List godoc
// @Summary Lists the projects you have access to
// @Description Lists the projects you have access to, this includes your own and the ones you were invited too not from orgs
// @Tags projects
// @ID projects_list
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fun.Response{data=[]models.Project}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /projects [get]
func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	namespaces, err := h.queries.ListProjects(r.Context())
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, namespaces)
}
