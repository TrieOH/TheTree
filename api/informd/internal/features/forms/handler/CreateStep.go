package handler

import (
	"Informd/models"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// CreateStep godoc
// @Summary Create a step
// @Description Creates a step on a form.
// @Tags steps
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param request body models.CreateStepRequest true "Form title"
// @Param form_id path string true "Form ID"
// @Success 201 {object} models.Step "Form created successfully"
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
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
	if bind.BailInto(w, req, &payload) {
		return
	}
	form, err := h.commands.CreateStep(r.Context(), formID, payload)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, form, http.StatusCreated)
}
