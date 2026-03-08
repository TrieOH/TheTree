package oauth_handler

import (
	"TriePayments/internal/core/application/oauth/commands"
	"TriePayments/internal/core/interfaces/http/dto"
	"TriePayments/internal/shared/validation"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

// BeginOAuth godoc
// @Summary Begin OAuth flow for a provider
// @Description Returns a redirect URL to start the OAuth authorization flow
// @Tags oauth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param provider path string true "Provider name (e.g. mercadopago)"
// @Param name path string true "Workspace name"
// @Param request body dto.BeginOAuthRequest true "OAuth request"
// @Success 200 {object} dto.BeginOAuthResponse
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /workspaces/{name}/oauth/{provider}/begin [post]
func (h *Handler) BeginOAuth(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	workspaceName := chi.URLParam(r, "name")

	var req dto.BeginOAuthRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	redirectURL, err := h.commands.BeginOAuth(r.Context(), commands.BeginOAuthRequest{
		Provider:         provider,
		WorkspaceName:    workspaceName,
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
