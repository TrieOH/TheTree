package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// RemoveMember godoc
// @Summary Remove an organization member
// @Description Lets you remove a member from the organization
// @Tags organizations
// @ID organizations_removemember
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.RemoveOrganizationMemberRequest true "Member details"
// @Success 201 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /organizations/{organization_id}/members [delete]
func (h *Handlers) RemoveMember(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	orgID, err := req.Path("organization_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.RemoveOrganizationMemberRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err = h.commands.RemoveMember(r.Context(), payload.ToInput(orgID))
	if fun.Bail(w, err) {
		return
	}
	fun.OK().Send(w)
}
