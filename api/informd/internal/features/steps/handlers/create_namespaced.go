package handlers

import (
	"net/http"

	"Informd/models"

	"github.com/MintzyG/fun"
)

// CreateNamespacedStep godoc
// @Summary creates a step in a namespaced form
// @Tags steps
// @ID steps_create_namespaced
// @Accept json
// @Produce json
// @Param namespace_id path string true "Namespace ID"
// @Param form_id path string true "Form ID"
// @Param request body models.CreateStepRequest true "Step payload"
// @Success 201 {object} fun.Response{data=models.Step}
// @Failure 400 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /namespaces/{namespace_id}/forms/{form_id}/steps [post]
func (h *Handlers) CreateNamespacedStep(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	namespaceID, err := req.Path("namespace_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	formID, err := req.Path("form_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.CreateStepRequest
	if fun.BailInto(w, req, &payload) {
		return
	}
	step, err := h.commands.CreateNamespaced(r.Context(), payload.ToNamespacedFormInput(namespaceID, formID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, step, http.StatusCreated)
}
