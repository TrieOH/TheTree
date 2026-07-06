package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// GetMemberByID godoc
// @Summary Get actors by ID
// @Tags organizations
// @ID organizations_getactorbyid
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fun.Response{data=[]idx.Actor}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Failure 503 {object} fun.Response
// @Router /organizations/{organization_id}/members/{member_id} [get]
func (h *Handlers) GetMemberByID(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	orgID, err := req.Path("organization_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	memberID, err := req.Path("member_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	members, err := h.queries.GetMemberByID(r.Context(), orgID, memberID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, members)
}
