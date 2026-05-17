package handler

import (
	"Informd/models"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
	"github.com/MintzyG/fun/middlewares"
)

// BulkGet godoc
// @Summary Bulk get forms
// @Description Returns a list of forms by their IDs. IDs should be obtained via a SpiceDB lookup on the client side.
// @Tags forms
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param request body models.BulkGetRequest true "Form IDs"
// @Success 200 {array} models.Form "Forms retrieved successfully"
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /forms/bulk [post]
func (h *Handler) BulkGet(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	params := middlewares.QueryParams[models.BulkGetParams](r)
	var payload models.BulkGetRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	forms, err := h.queries.BulkGet(r.Context(), payload.IDs, params)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, forms)
}
