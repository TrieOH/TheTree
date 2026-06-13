package handlers

import (
	"net/http"

	"Informd/models"
	"lib/xslices"

	"github.com/MintzyG/fun"
)

// BulkEditSteps godoc
// @Summary bulk edits steps in a form
// @Tags steps
// @ID steps_bulk_edit
// @Accept json
// @Produce json
// @Param form_id path string true "Form ID"
// @Param request body []models.UpdateStepRequest true "Steps payload"
// @Success 200 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms/{form_id}/steps [put]
func (h *Handlers) BulkEditSteps(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload []models.UpdateStepRequest
	if fun.BailInto(w, req, &payload) {
		return
	}
	inputs := xslices.MapSlice(payload, func(s models.UpdateStepRequest) models.UpdateFormStepInput {
		return s.ToFormInput(formID)
	})
	err = h.commands.BulkEdit(r.Context(), formID, inputs)
	if fun.Bail(w, err) {
		return
	}
	fun.OK().Send(w)
}
