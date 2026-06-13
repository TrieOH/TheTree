package handlers

import (
	"net/http"

	"Informd/models"

	"github.com/MintzyG/fun"
)

// EditSelectConfig godoc
// @Summary edits the select config for a field
// @Tags fields
// @ID fields_edit_select_config
// @Accept json
// @Produce json
// @Param form_id path string true "Form ID"
// @Param step_id path string true "Step ID"
// @Param field_id path string true "Field ID"
// @Param request body models.FieldSelectConfig true "Select config payload"
// @Success 200 {object} fun.Response{data=models.FieldSelectConfig}
// @Failure 400 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms/{form_id}/steps/{step_id}/fields/{field_id}/select [put]
func (h *Handlers) EditSelectConfig(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	fieldID, err := req.Path("field_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.FieldSelectConfig
	if fun.BailInto(w, req, &payload) {
		return
	}
	payload.FieldID = fieldID
	config, err := h.commands.EditSelectConfig(r.Context(), formID, payload)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, config)
}
