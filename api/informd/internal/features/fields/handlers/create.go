package handlers

import (
	"net/http"

	"Informd/models"

	"github.com/MintzyG/fun"
)

// CreateField godoc
// @Summary creates a field in a step
// @Tags fields
// @ID fields_create
// @Accept json
// @Produce json
// @Param form_id path string true "Form ID"
// @Param step_id path string true "Step ID"
// @Param request body models.CreateFieldRequest true "Field payload"
// @Success 201 {object} fun.Response{data=models.Field}
// @Failure 400 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms/{form_id}/steps/{step_id}/fields [post]
func (h *Handlers) CreateField(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	stepID, err := req.Path("step_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.CreateFieldRequest
	if fun.BailInto(w, req, &payload) {
		return
	}
	field, err := h.commands.Create(r.Context(), payload.ToStepInput(formID, stepID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, field, http.StatusCreated)
}
