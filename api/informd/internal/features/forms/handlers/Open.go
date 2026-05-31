package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// Open godoc
// @Summary Open a draft form for responses
// @Tags forms
// @ID forms_open
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param form_id path string true "Form ID"
// @Success 201 {object} models.Form "Form updated successfully"
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms/{form_id}/open [post]
func (h *Handlers) Open(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	form, err := h.commands.Open(r.Context(), formID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, form)
}
