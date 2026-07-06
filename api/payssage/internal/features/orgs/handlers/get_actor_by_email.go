package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// GetMemberByEmail godoc
// @Summary Get actors by Email
// @Tags organizations
// @ID organizations_getactorbyemail
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fun.Response{data=[]idx.Actor}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /organizations/{organization_id}/members/{member_email}:by_email [get]
func (h *Handlers) GetMemberByEmail(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	orgID, err := req.Path("organization_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	memberEmail, err := req.Path("member_email").StringRequired()
	if fun.Bail(w, err) {
		return
	}
	members, err := h.queries.GetMemberByEmail(r.Context(), memberEmail, orgID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, members)
}
