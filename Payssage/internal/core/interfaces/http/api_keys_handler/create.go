package workspaces_handler

import (
	"TriePayments/internal/core/interfaces/http/dto"
	"TriePayments/internal/shared/validation"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

// Create godoc
// @Summary Create an API key
// @Description Creates a new API key for the given workspace. The raw key is returned once and never stored.
// @Tags workspaces
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Param request body dto.CreateAPIKeyRequest true "API key details"
// @Success 201 {object} dto.CreateAPIKeyResponse "API key created successfully"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /workspaces/{name}/keys [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")

	var req dto.CreateAPIKeyRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	rawKey, apiKey, err := h.commands.Create(r.Context(), workspaceName, req.Name)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created().WithData(dto.CreateAPIKeyResponse{
		APIKeyResponse: dto.APIKeyResponse{
			ID:        apiKey.ID,
			Name:      apiKey.Name,
			Prefix:    apiKey.KeyPrefix,
			CreatedAt: apiKey.CreatedAt,
			RevokedAt: apiKey.RevokedAt,
		},
		Key: rawKey,
	}).Send(w)
}
