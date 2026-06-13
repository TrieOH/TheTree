package handlers

import (
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
)

// ListMembers godoc
// @Summary Lists the organization members
// @Description Gets a list of all the members added to the organization, includes the owner
// @Tags organizations
// @ID organizations_listmembers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fun.Response{data=[]models.OrganizationMember}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /organizations/{organization_id}/members [get]
func (h *Handlers) ListMembers(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	orgID, err := req.Path("organization_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	members, err := h.queries.ListMembers(r.Context(), orgID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, members)
}
