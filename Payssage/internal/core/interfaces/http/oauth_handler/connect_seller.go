package oauth_handler

import (
	"TriePayments/internal/core/application/oauth/commands"
	"TriePayments/internal/core/interfaces/http/dto"
	"TriePayments/internal/shared/validation"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

// ConnectSeller godoc
// @Summary Connect a seller account to a workspace
// @Description Begins OAuth flow for a seller to connect their account for split payments
// @Tags oauth
// @Accept json
// @Produce json
// @Param X-API-Key header string true "X-API-Key: tp_xxxxxxxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Param provider path string true "Provider name (e.g. mercadopago)"
// @Param request body dto.ConnectSellerRequest true "Connect request"
// @Success 200 {object} dto.BeginOAuthResponse
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /workspaces/{name}/providers/{provider}/connect [post]
func (h *Handler) ConnectSeller(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")
	provider := chi.URLParam(r, "provider")

	var req dto.ConnectSellerRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	redirectURL, finalRedirectURL, err := h.commands.ConnectSeller(r.Context(), commands.ConnectSellerRequest{
		WorkspaceName:       workspaceName,
		Provider:            provider,
		ProviderRedirectURL: req.ProviderRedirectURL,
		FinalRedirectURL:    req.FinalRedirectURL,
	})
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.BeginOAuthResponse{
		RedirectURL:      redirectURL,
		FinalRedirectURL: finalRedirectURL,
	}).Send(w)
}
