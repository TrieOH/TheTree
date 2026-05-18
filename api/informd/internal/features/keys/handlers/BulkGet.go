package handlers

import (
	"Informd/models"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// BulkGet godoc
// @Summary Bulk get api keys
// @Description Returns a list of api keys by their IDs. IDs should be obtained via a SpiceDB lookup on the client side.
// @Tags api_keys
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param request body models.BulkGetRequest true "APIKey IDs"
// @Success 200 {array} models.Form "Forms retrieved successfully"
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /api-keys/bulk [post]
func (h *Handler) BulkGet(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	var payload models.BulkGetRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	forms, err := h.queries.BulkGet(r.Context(), payload.IDs)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, forms)
}
