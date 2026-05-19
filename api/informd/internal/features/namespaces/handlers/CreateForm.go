package handlers

import (
	"Informd/models"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// CreateForm godoc
// @Summary Create a form
// @Description Creates a namespaced form.
// @Tags namespaces
// @ID namespaces_createform
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateFormRequest true "Form title"
// @Success 201 {object} fun.Response{data=models.Form} "Form created successfully"
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /namespaces/{namespace_id}/forms [post]
func (h *Handlers) CreateForm(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	namespaceID, err := req.Path("namespace_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.CreateFormRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	form, err := h.formsCommands.Create(r.Context(), payload.Title, &namespaceID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, form, http.StatusCreated)
}
