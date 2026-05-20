package handlers

import (
	"Informd/models"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// Create godoc
// @Summary Create a form
// @Description Creates a form not namespaced.
// @Tags forms
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param request body models.CreateFormRequest true "Form title"
// @Success 201 {object} models.Form "Form created successfully"
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms [post]
func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	var payload models.CreateFormRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	form, err := h.commands.Create(r.Context(), payload.Title)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, form, http.StatusCreated)
}
