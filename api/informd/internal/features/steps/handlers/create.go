package handlers

import (
	"net/http"

	"Informd/models"

	"github.com/MintzyG/fun"
)

// CreateStep godoc
// @Summary creates a step in a form
// @Tags steps
// @ID steps_create
// @Accept json
// @Produce json
// @Param form_id path string true "Form ID"
// @Param request body models.CreateStepRequest true "Step payload"
// @Success 201 {object} fun.Response{data=models.Step}
// @Failure 400 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms/{form_id}/steps [post]
func (h *Handlers) CreateStep(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.CreateStepRequest
	if fun.BailInto(w, req, &payload) {
		return
	}
	step, err := h.commands.Create(r.Context(), payload.ToFormInput(formID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, step, http.StatusCreated)
}
