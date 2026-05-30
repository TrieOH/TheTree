package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// ListNamespaced godoc
// @Summary lists fields in a namespaced step
// @Tags fields
// @ID fields_list_namespaced
// @Accept json
// @Produce json
// @Param namespace_id path string true "Namespace ID"
// @Param form_id path string true "Form ID"
// @Param step_id path string true "Step ID"
// @Success 200 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields [get]
func (h *Handlers) ListNamespaced(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	namespaceID, err := req.Path("namespace_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	stepID, err := req.Path("step_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	fields, err := h.queries.ListNamespaced(r.Context(), formID, namespaceID, stepID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, fields)
}
