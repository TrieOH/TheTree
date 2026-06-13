package handlers

import (
	"net/http"

	"Informd/models"

	"github.com/MintzyG/fun"
)

// EditSelectConfigNamespaced godoc
// @Summary edits the select config for a field in a namespaced step
// @Tags fields
// @ID fields_edit_select_config_namespaced
// @Accept json
// @Produce json
// @Param namespace_id path string true "Namespace ID"
// @Param form_id path string true "Form ID"
// @Param step_id path string true "Step ID"
// @Param field_id path string true "Field ID"
// @Param request body models.FieldSelectConfig true "Select config payload"
// @Success 200 {object} fun.Response{data=models.FieldSelectConfig}
// @Failure 400 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields/{field_id}/select [put]
func (h *Handlers) EditSelectConfigNamespaced(w http.ResponseWriter, r *http.Request) {
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
	var payload models.FieldSelectConfig
	if fun.BailInto(w, req, &payload) {
		return
	}
	payload.FieldID = fieldID
	config, err := h.commands.EditSelectConfigNamespaced(r.Context(), formID, namespaceID, payload)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, config)
}
