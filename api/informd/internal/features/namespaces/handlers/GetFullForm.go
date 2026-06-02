package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// GetFullFormNamespaced godoc
// @Summary Get full namespaced form with responses
// @Description Gets the full namespaced form structure including steps, fields, and all submitted answers with responder info
// @Tags forms
// @ID forms_getfullform
// @Produce json
// @Security BearerAuth
// @Param namespace_id path string true "Namespace ID"
// @Param form_id path string true "Form ID"
// @Success 200 {object} fun.Response{data=models.FullForm}
// @Failure 401 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /namespaces/{namespace_id}/forms/{form_id}/full [get]
func (h *Handlers) GetFullFormNamespaced(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	namespaceID, err := req.Path("namespace_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	form, err := h.queries.GetFullForm(r.Context(), namespaceID, formID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, form)
}
