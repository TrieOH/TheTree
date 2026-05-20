package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// ListFormMembers godoc
// @Summary Lists the members of a form
// @Description Gets a list of all the members added to the namespaced form
// @Tags namespaces
// @ID namespaces_listformmembers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fun.Response{data=[]models.FormMember}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /namespaces/{namespace_id}/forms/{form_id}/members [get]
func (h *Handlers) ListFormMembers(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	namespaceID, err := req.Path("namespace_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	members, err := h.queries.ListFormMembers(r.Context(), namespaceID, formID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, members)
}
