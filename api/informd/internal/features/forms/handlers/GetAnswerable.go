package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// GetAnswerable godoc
// @Summary Get full form with responses for answering
// @Description Gets the full form structure including steps, fields for asnwering
// @Tags forms
// @ID forms_getanswerable
// @Produce json
// @Security BearerAuth
// @Param form_id path string true "Form ID"
// @Success 200 {object} fun.Response{data=models.FormAnswerable}
// @Failure 401 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms/{form_id}/asnswerable [get]
func (h *Handlers) GetAnswerable(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	form, err := h.queries.GetAnswerable(r.Context(), formID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, form)
}
