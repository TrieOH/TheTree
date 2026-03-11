package oauth_handler

import (
	"TriePayments/internal/core/interfaces/http/dto"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

// ListMarketplaceConfigs godoc
// @Summary List marketplace configurations for a workspace
// @Description Returns all marketplace provider configs for the workspace
// @Tags oauth
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Success 200 {object} []dto.MarketplaceConfigResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /workspaces/{name}/marketplace [get]
func (h *Handler) ListMarketplaceConfigs(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")

	configs, err := h.queries.ListMarketplaceConfigs(r.Context(), workspaceName)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	result := make([]dto.MarketplaceConfigResponse, len(configs))
	for i, config := range configs {
		result[i] = dto.MarketplaceConfigResponse{
			ID:           config.ID,
			WorkspaceID:  config.WorkspaceID,
			CredentialID: config.CredentialID,
			Provider:     config.Provider,
			FeeBps:       config.FeeBps,
			CreatedAt:    config.CreatedAt,
			UpdatedAt:    config.UpdatedAt,
		}
	}

	resp.OK().WithData(result).Send(w)
}
