package oauth_handler

import (
	"TriePayments/internal/core/application/oauth/commands"
	"TriePayments/internal/core/interfaces/http/dto"
	"TriePayments/internal/shared/validation"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

// SetMarketplaceConfig godoc
// @Summary Configure marketplace settings for a workspace
// @Description Sets the MP credential and platform fee for marketplace split payments
// @Tags oauth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Param request body dto.SetMarketplaceConfigRequest true "Marketplace config"
// @Success 200 {object} dto.MarketplaceConfigResponse
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 403 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /workspaces/{name}/marketplace [put]
func (h *Handler) SetMarketplaceConfig(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")

	var req dto.SetMarketplaceConfigRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	config, err := h.commands.SetMarketplaceConfig(r.Context(), commands.SetMarketplaceConfigRequest{
		WorkspaceName: workspaceName,
		CredentialID:  req.CredentialID,
		FeeBps:        req.FeeBps,
	})
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.MarketplaceConfigResponse{
		ID:           config.ID,
		WorkspaceID:  config.WorkspaceID,
		CredentialID: config.CredentialID,
		Provider:     config.Provider,
		FeeBps:       config.FeeBps,
		CreatedAt:    config.CreatedAt,
		UpdatedAt:    config.UpdatedAt,
	}).Send(w)
}
