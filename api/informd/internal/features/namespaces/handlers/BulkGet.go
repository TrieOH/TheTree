package handlers

import (
	"Informd/models"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// BulkGet godoc
// @Summary Bulk get namespaces
// @Description Returns a list of namespaces by their IDs. IDs should be obtained via a SpiceDB lookup on the client side.
// @Tags namespaces
// @ID namespaces_bulkget
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.BulkGetRequest true "Namespace IDs"
// @Success 200 {array} models.Form "Namespaces retrieved successfully"
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /namespaces/bulk [post]
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
