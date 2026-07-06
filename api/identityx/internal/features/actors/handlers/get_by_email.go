package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

// GetByEmail godoc
// @Summary Get actors by Email
// @Tags organizations
// @ID organizations_getactorbyemail
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fun.Response{data=models.Actor}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /projects/{project_id}/actors/{actor_email}:by_email [get]
func (h *Handlers) GetByEmail(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	projectID, err := req.Path("project_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	email, err := req.Path("actor_email").StringRequired()
	if fun.Bail(w, err) {
		return
	}
	members, err := h.queries.GetByEmail(r.Context(), email, projectID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, members)
}
