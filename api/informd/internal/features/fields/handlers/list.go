package handlers

import (
	"net/http"

	"github.com/MintzyG/fun"
)

// List godoc
// @Summary lists fields in a step
// @Tags fields
// @ID fields_list
// @Accept json
// @Produce json
// @Param form_id path string true "Form ID"
// @Param step_id path string true "Step ID"
// @Success 200 {object} fun.Response
// @Failure 400 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms/{form_id}/steps/{step_id}/fields [get]
func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	stepID, err := req.Path("step_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	fields, err := h.queries.List(r.Context(), formID, stepID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, fields)
}
