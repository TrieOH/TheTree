package handlers

import (
	"net/http"

	"Informd/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// AddFormMember godoc
// @Summary Add a form member
// @Description Lets you add a member to the namespaced form as a viewer, editor or admin
// @Tags namespaces
// @ID namespaces_addformmember
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.AddFormMemberRequest true "Member details"
// @Success 201 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /namespaces/{namespace_id}/forms/{form_id}/members [post]
func (h *Handlers) AddFormMember(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	namespaceID, err := req.Path("namespace_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.AddFormMemberRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err = h.commands.AddFormMember(r.Context(), payload.ToNamespaceInput(namespaceID, formID))
	if fun.Bail(w, err) {
		return
	}
	fun.Created().Send(w)
}
