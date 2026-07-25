package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

func (h *Handlers) RemoveProjectMember(w http.ResponseWriter, r *http.Request) {
	if !globals.SetupComplete() {
		fun.ServiceUnavailable("please setup IDX first on /auth/setup").Send(w)
		return
	}
	req := fun.From(r)
	orgID, err := req.Path("organization_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	projectID, err := req.Path("project_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.RemoveOrgProjectMemberRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err = h.commands.RemoveProjectMember(r.Context(), payload.ToInput(orgID, projectID))
	if fun.Bail(w, err) {
		return
	}
	fun.OK().Send(w)
}
