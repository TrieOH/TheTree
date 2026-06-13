package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// AddMember godoc
// @Summary Add a organization member
// @Description Lets you add a member to the organization as a member or admin
// @Tags organizations
// @ID organizations_addmember
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.AddOrganizationMemberRequest true "Member details"
// @Success 201 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /organizations/{organization_id}/members [post]
func (h *Handlers) AddMember(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	orgID, err := req.Path("organization_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.AddOrganizationMemberRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err = h.commands.AddMember(r.Context(), payload.ToInput(orgID))
	if fun.Bail(w, err) {
		return
	}
	fun.Created().Send(w)
}
