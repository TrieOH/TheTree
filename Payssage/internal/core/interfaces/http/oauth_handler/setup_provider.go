package oauth_handler

import (
	"TriePayments/internal/core/application/oauth/commands"
	"TriePayments/internal/core/interfaces/http/dto"
	"TriePayments/internal/shared/validation"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

// SetupProvider godoc
// @Summary Set up a payment provider for a workspace
// @Description Begins OAuth flow to connect a payment provider to the workspace
// @Tags oauth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Param provider path string true "Provider name (e.g. mercadopago)"
// @Param request body dto.SetupProviderRequest true "Setup request"
// @Success 200 {object} dto.BeginOAuthResponse
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /workspaces/{name}/providers/{provider}/setup [post]
func (h *Handler) SetupProvider(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")
	provider := chi.URLParam(r, "provider")

	var req dto.SetupProviderRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	redirectURL, err := h.commands.SetupProvider(r.Context(), commands.SetupProviderRequest{
		WorkspaceName:    workspaceName,
		Provider:         provider,
		IsMarketplace:    req.IsMarketplace,
		FeeBps:           req.FeeBps,
		FinalRedirectURL: req.FinalRedirectURL,
	})
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.BeginOAuthResponse{
		RedirectURL: redirectURL,
	}).Send(w)
}
