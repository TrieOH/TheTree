package handlers

import (
	"IdentityX/models"
	"lib/globals"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

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
