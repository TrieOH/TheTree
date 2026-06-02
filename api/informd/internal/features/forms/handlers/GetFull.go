package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// GetFull godoc
// @Summary Get full form with responses
// @Description Gets the full form structure including steps, fields, and all submitted answers with responder info
// @Tags forms
// @ID forms_getfull
// @Produce json
// @Security BearerAuth
// @Param form_id path string true "Form ID"
// @Success 200 {object} fun.Response{data=models.FullForm}
// @Failure 401 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms/{form_id}/full [get]
func (h *Handlers) GetFull(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	form, err := h.queries.GetFull(r.Context(), formID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, form)
}
