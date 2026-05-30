package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// DeleteField godoc
// @Summary deletes a field from a step
// @Tags fields
// @ID fields_delete
// @Accept json
// @Produce json
// @Param form_id path string true "Form ID"
// @Param step_id path string true "Step ID"
// @Param field_id path string true "Field ID"
// @Success 200 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms/{form_id}/steps/{step_id}/fields/{field_id} [delete]
func (h *Handlers) DeleteField(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	fieldID, err := req.Path("field_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	err = h.commands.Delete(r.Context(), formID, fieldID)
	if fun.Bail(w, err) {
		return
	}
	fun.OK().Send(w)
}
