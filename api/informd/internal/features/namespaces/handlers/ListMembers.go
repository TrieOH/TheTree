package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// ListMembers godoc
// @Summary Lists the namespace members
// @Description Gets a list of all the members added to the namespace, includes the owner
// @Tags namespaces
// @ID namespaces_listmembers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} fun.Response{data=[]models.NamespaceMember}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /namespaces/{namespace_id}/members [get]
func (h *Handler) ListMembers(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	namespaceID, err := req.Path("namespace_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	members, err := h.queries.ListMembers(r.Context(), namespaceID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, members)
}
