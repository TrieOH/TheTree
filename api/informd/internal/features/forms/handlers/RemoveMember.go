package handlers

import (
	"net/http"

	"Informd/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// RemoveMember godoc
// @Summary Remove a form member
// @Description Lets you remove a member from the form
// @Tags forms
// @ID forms_removemember
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param form_id path string true "Form ID"
// @Param request body models.RemoveFormMemberRequest true "Member details"
// @Success 201 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms/{form_id}/members [delete]
func (h *Handlers) RemoveMember(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.RemoveFormMemberRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err = h.commands.RemoveMember(r.Context(), payload.ToInput(formID))
	if fun.Bail(w, err) {
		return
	}
	fun.Created().Send(w)
}
