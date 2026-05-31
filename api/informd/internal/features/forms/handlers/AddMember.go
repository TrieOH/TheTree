package handlers

import (
	"Informd/models"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// AddMember godoc
// @Summary Add a form member
// @Description Lets you add a member to the form as a viewer, editor or admin
// @Tags forms
// @ID forms_addmember
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param form_id path string true "Form ID"
// @Param request body models.AddFormMemberRequest true "Member details"
// @Success 201 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms/{form_id}/members [post]
func (h *Handlers) AddMember(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.AddFormMemberRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err = h.commands.AddMember(r.Context(), payload.ToInput(formID))
	if fun.Bail(w, err) {
		return
	}
	fun.Created().Send(w)
}
