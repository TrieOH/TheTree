package handlers

import (
	"net/http"

	"Informd/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// RemoveMember godoc
// @Summary Remove a namespace member
// @Description Lets you remove a member from the namespace
// @Tags namespaces
// @ID namespaces_removemember
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.RemoveNamespaceMemberRequest true "Member details"
// @Success 201 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /namespaces/{namespace_id}/members [delete]
func (h *Handlers) RemoveMember(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	namespaceID, err := req.Path("namespace_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.RemoveNamespaceMemberRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err = h.commands.RemoveMember(r.Context(), payload.ToInput(namespaceID))
	if fun.Bail(w, err) {
		return
	}
	fun.OK().Send(w)
}
