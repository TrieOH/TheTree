package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// ListMembers godoc
// @Summary Lists the form members
// @Description Gets a list of all the members added to the form, includes the owner
// @Tags forms
// @ID forms_listmembers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fun.Response{data=[]models.FormMember}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms/{form_id}/members [get]
func (h *Handlers) ListMembers(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	members, err := h.queries.ListMembers(r.Context(), formID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, members)
}
