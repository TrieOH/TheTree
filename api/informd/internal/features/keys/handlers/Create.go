package handlers

import (
	"Informd/models"
	"net/http"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
)

// Create godoc
// @Summary Create an API key
// @Description Creates a new API key for the given project. The raw key is returned once and never stored.
// @Tags api_keys
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param project_id path string true "Project ID"
// @Param request body models.CreateAPIKeyRequest true "API key details"
// @Success 201 {object} models.CreateAPIKeyResponse "API key created successfully"
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /api-keys [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	var payload models.CreateAPIKeyRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	rawKey, apiKey, err := h.commands.Create(r.Context(), payload.Name)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, models.CreateAPIKeyResponse{
		APIKeyResponse: apiKey.ToResponse(),
		Key:            rawKey,
	}, http.StatusCreated)
}
