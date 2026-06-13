package handlers

import (
	"net/http"

	"Informd/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// RemoveFormMember godoc
// @Summary Remove a form member
// @Description Lets you remove a member from the namespaced form
// @Tags namespaces
// @ID namespaces_removeformmember
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.RemoveFormMemberRequest true "Member details"
// @Success 201 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /namespaces/{namespace_id}/forms/{form_id}/members [delete]
func (h *Handlers) RemoveFormMember(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	namespaceID, err := req.Path("namespace_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.RemoveFormMemberRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err = h.commands.RemoveFormMember(r.Context(), payload.ToNamespaceInput(namespaceID, formID))
	if fun.Bail(w, err) {
		return
	}
	fun.OK().Send(w)
}
