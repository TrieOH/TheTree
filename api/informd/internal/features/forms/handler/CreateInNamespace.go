package handler

import (
	"Informd/models"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// CreateInNamespace godoc
// @Summary Create a form
// @Description Creates a form in the given namespace.
// @Tags forms
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param namespace_id path string true "Namespace ID"
// @Param request body models.CreateFormRequest true "Form title"
// @Success 201 {object} models.Form "Form created successfully"
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /namespaces/{namespace_id}/forms [post]
func (h *Handler) CreateInNamespace(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	namespaceID, err := req.Path("namespace_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload models.CreateFormRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	form, err := h.commands.Create(r.Context(), payload.Title, &namespaceID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, form, http.StatusCreated)
}
