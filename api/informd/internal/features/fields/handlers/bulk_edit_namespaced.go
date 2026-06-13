package handlers

import (
	"net/http"

	"Informd/models"
	"lib/xslices"

	"github.com/MintzyG/fun"
)

// BulkEditNamespacedFields godoc
// @Summary bulk edits fields in a namespaced step
// @Tags fields
// @ID fields_bulk_edit_namespaced
// @Accept json
// @Produce json
// @Param namespace_id path string true "Namespace ID"
// @Param form_id path string true "Form ID"
// @Param step_id path string true "Step ID"
// @Param request body []models.UpdateFieldRequest true "Fields payload"
// @Success 200 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields [put]
func (h *Handlers) BulkEditNamespacedFields(w http.ResponseWriter, r *http.Request) {
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
	var payload []models.UpdateFieldRequest
	if fun.BailInto(w, req, &payload) {
		return
	}
	inputs := xslices.MapSlice(payload, func(f models.UpdateFieldRequest) models.UpdateNamespacedStepFieldInput {
		return f.ToNamespacedStepInput(namespaceID, formID, stepID)
	})
	err = h.commands.BulkEditNamespaced(r.Context(), formID, namespaceID, inputs)
	if fun.Bail(w, err) {
		return
	}
	fun.OK().Send(w)
}
