package handlers

import (
	"Informd/models"
	"lib/xslices"
	"net/http"

	"github.com/MintzyG/fun"
)

// BulkEditFields godoc
// @Summary bulk edits fields in a step
// @Tags fields
// @ID fields_bulk_edit
// @Accept json
// @Produce json
// @Param form_id path string true "Form ID"
// @Param step_id path string true "Step ID"
// @Param request body []models.UpdateFieldRequest true "Fields payload"
// @Success 200 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms/{form_id}/steps/{step_id}/fields [put]
func (h *Handlers) BulkEditFields(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
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
	inputs := xslices.MapSlice(payload, func(f models.UpdateFieldRequest) models.UpdateStepFieldInput {
		return f.ToStepInput(stepID)
	})
	err = h.commands.BulkEdit(r.Context(), formID, inputs)
	if fun.Bail(w, err) {
		return
	}
	fun.OK().Send(w)
}
