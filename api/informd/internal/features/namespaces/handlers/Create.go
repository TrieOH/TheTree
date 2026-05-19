package handlers

import (
	"Informd/models"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// Create godoc
// @Summary Create a namespace
// @Description Creates a new namespace for the authenticated user
// @Tags namespaces
// @ID namespaces_create
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateNamespaceRequest true "Project details"
// @Success 201 {object} fun.Response{data=models.Namespace}
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /namespaces [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	var payload models.CreateNamespaceRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	namespace, err := h.commands.Create(r.Context(), payload.Name)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, namespace, http.StatusCreated)
}
