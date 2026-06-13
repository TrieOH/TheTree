package handlers

import (
	"net/http"

	"Informd/models"

	"github.com/MintzyG/fun"
)

// CreateNamespacedField godoc
// @Summary creates a field in a namespaced step
// @Tags fields
// @ID fields_create_namespaced
// @Accept json
// @Produce json
// @Param namespace_id path string true "Namespace ID"
// @Param form_id path string true "Form ID"
// @Param step_id path string true "Step ID"
// @Param request body models.CreateFieldRequest true "Field payload"
// @Success 201 {object} fun.Response{data=models.Field}
// @Failure 400 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /namespaces/{namespace_id}/forms/{form_id}/steps/{step_id}/fields [post]
func (h *Handlers) CreateNamespacedField(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	namespaceID, err := req.Path("namespace_id").UUID()
	if fun.Bail(w, err) {
		return
	}
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
	field, err := h.commands.CreateNamespaced(r.Context(), payload.ToNamespacedStepInput(namespaceID, formID, stepID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, field, http.StatusCreated)
}
