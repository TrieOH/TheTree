package handlers

import (
	"net/http"

	"Informd/models"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// Submit godoc
// @Summary Submit a form response
// @Tags forms
// @ID forms_submit
// @Accept json
// @Param request body models.SubmitRequest true "Form title"
// @Success 201 {object} fun.Response "Form answered successfully"
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms/{form_id}/responses [post]
func (h *Handlers) Submit(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.SubmitRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	err = h.commands.Submit(r.Context(), payload.ToInput(formID))
	if fun.Bail(w, err) {
		return
	}
	fun.Created().Send(w)
}
