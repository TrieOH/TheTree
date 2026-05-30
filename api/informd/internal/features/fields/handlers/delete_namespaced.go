package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// DeleteNamespacedField godoc
// @Summary deletes a field from a namespaced step
// @Tags fields
// @ID fields_delete_namespaced
// @Accept json
// @Produce json
// @Param namespace_id path string true "Namespace ID"
// @Param form_id path string true "Form ID"
// @Param step_id path string true "Step ID"
// @Param field_id path string true "Field ID"
// @Success 200 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields/{field_id} [delete]
func (h *Handlers) DeleteNamespacedField(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	namespaceID, err := req.Path("namespace_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	fieldID, err := req.Path("field_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	err = h.commands.DeleteNamespaced(r.Context(), namespaceID, formID, fieldID)
	if fun.Bail(w, err) {
		return
	}
	fun.OK().Send(w)
}
