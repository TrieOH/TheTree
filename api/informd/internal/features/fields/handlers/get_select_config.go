package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// GetSelectConfig godoc
// @Summary gets the select config for a field
// @Tags fields
// @ID fields_get_select_config
// @Accept json
// @Produce json
// @Param form_id path string true "Form ID"
// @Param step_id path string true "Step ID"
// @Param field_id path string true "Field ID"
// @Success 200 {object} fun.Response{data=models.FieldSelectConfig}
// @Failure 400 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms/{form_id}/steps/{step_id}/fields/{field_id}/select [get]
func (h *Handlers) GetSelectConfig(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	fieldID, err := req.Path("field_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	config, err := h.queries.GetSelectConfig(r.Context(), formID, fieldID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, config)
}
