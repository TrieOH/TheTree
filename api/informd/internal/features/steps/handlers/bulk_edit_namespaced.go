package handlers

import (
	"net/http"

	"Informd/models"
	"lib/xslices"

	"github.com/MintzyG/fun"
)

// BulkEditNamespacedSteps godoc
// @Summary bulk edits steps in a namespaced form
// @Tags steps
// @ID steps_bulk_edit_namespaced
// @Accept json
// @Produce json
// @Param namespace_id path string true "Namespace ID"
// @Param form_id path string true "Form ID"
// @Param request body []models.UpdateStepRequest true "Steps payload"
// @Success 200 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /namespaces/{namespace_id}/forms/{form_id}/steps [put]
func (h *Handlers) BulkEditNamespacedSteps(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	namespaceID, err := req.Path("namespace_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload []models.UpdateStepRequest
	if fun.BailInto(w, req, &payload) {
		return
	}
	inputs := xslices.MapSlice(payload, func(s models.UpdateStepRequest) models.UpdateNamespacedFormStepInput {
		return s.ToNamespacedFormInput(namespaceID, formID)
	})
	err = h.commands.BulkEditNamespaced(r.Context(), formID, namespaceID, inputs)
	if fun.Bail(w, err) {
		return
	}
	fun.OK().Send(w)
}
