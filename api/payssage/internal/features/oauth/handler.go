package oauth

import (
	"lib/telemetry"
	"net/http"

	"payssage/internal/shared/errx"
	"payssage/internal/shared/validation"

	_ "payssage/models"

	"github.com/MintzyG/fun"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Handler struct {
	commands *CommandService
	queries  *QueryService
}

func NewHandler(
	commands *CommandService,
	queries *QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

// CompleteOAuth godoc
// @Summary OAuth callback from provider
// @Description Handles the provider callback, exchanges code for token, stores credential, redirects to final URL
// @Tags oauth
// @Param provider path string true "Provider name (e.g. mercadopago)"
// @Param code query string true "Authorization code from provider"
// @Param state query string true "State token"
// @Param redirect_uri query string true "Redirect URI"
// @Failure 400 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /oauth/{provider}/callback [get]
func (h *Handler) CompleteOAuth(w http.ResponseWriter, r *http.Request) {
	provider := chi.URLParam(r, "provider")
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	redirectURI := r.URL.Query().Get("redirect_uri")

	if code == "" || state == "" {
		fun.BadRequest("code and state are required").Send(w)
		return
	}

	finalURL, err := h.commands.CompleteOAuth(r.Context(), provider, state, code, redirectURI)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(map[string]string{
		"url": finalURL,
	}).Send(w)
}

type ConnectSellerRequest struct {
	ProviderRedirectURL string `json:"provider_redirect_url" validate:"required,url"`
	FinalRedirectURL    string `json:"final_redirect_url" validate:"required,url"`
}

type BeginOAuthResponse struct {
	RedirectURL      string `json:"redirect_url"`
	FinalRedirectURL string `json:"final_redirect_url"`
}

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
// @Param request body ConnectSellerRequest true "Connect request"
// @Success 200 {object} BeginOAuthResponse
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /workspaces/{name}/providers/{provider}/connect [post]
func (h *Handler) ConnectSeller(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")
	provider := chi.URLParam(r, "provider")

	var req ConnectSellerRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		fun.Error(err).Send(w)
		return
	}

	redirectURL, finalRedirectURL, err := h.commands.ConnectSeller(r.Context(), ConnectSellerInput{
		WorkspaceName:       workspaceName,
		Provider:            provider,
		ProviderRedirectURL: req.ProviderRedirectURL,
		FinalRedirectURL:    req.FinalRedirectURL,
	})
	if err != nil {
		telemetry.Log().Info("Error connecting seller", zap.Error(err))
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(BeginOAuthResponse{
		RedirectURL:      redirectURL,
		FinalRedirectURL: finalRedirectURL,
	}).Send(w)
}

// DeleteMarketplaceConfig godoc
// @Summary Remove marketplace configuration
// @Description Removes the marketplace config for a workspace, reverting to simple mode
// @Tags oauth
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Success 200 {object} object
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /workspaces/{name}/marketplace/{credential_id} [delete]
func (h *Handler) DeleteMarketplaceConfig(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")
	credentialIDStr := chi.URLParam(r, "credential_id")

	credentialID, err := uuid.Parse(credentialIDStr)
	if err != nil {
		fun.BadRequest("invalid credential_id").Send(w)
		return
	}

	if err := h.commands.DeleteMarketplaceConfig(r.Context(), workspaceName, credentialID); err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().Send(w)
}

// DisconnectProvider godoc
// @Summary Disconnect a provider credential (seller)
// @Description Called via API key by Univents when a seller clicks Disconnect
// @Tags providers
// @Param X-API-Key header string true "X-API-Key: tp_xxxxxxxx"
// @Security APIKey
// @Param name path string true "Workspace name"
// @Param credential_id path string true "Credential ID"
// @Success 200 {object} object "disconnected successfully"
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /workspaces/{name}/providers/{credential_id}/disconnect [delete]
func (h *Handler) DisconnectProvider(w http.ResponseWriter, r *http.Request) {
	credentialID, rs := validation.GetUUID(r, "credential_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	_, err := h.commands.DisconnectCredential(r.Context(), credentialID)
	if err != nil {
		if errx.IsKind(err, "not_found") {
			fun.NotFound("credential not found").Send(w)
			return
		}
		fun.Error(err).Send(w)
		return
	}

	fun.OK("disconnected successfully").Send(w)
}

// ListMarketplaceConfigs godoc
// @Summary List marketplace configurations for a workspace
// @Description Returns all marketplace provider configs for the workspace
// @Tags oauth
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Success 200 {array} models.MarketplaceConfig
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /workspaces/{name}/marketplace [get]
func (h *Handler) ListMarketplaceConfigs(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")

	configs, err := h.queries.ListMarketplaceConfigs(r.Context(), workspaceName)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(configs).Send(w)
}

// RevokeProvider godoc
// @Summary Revoke a provider credential (owner)
// @Description Workspace owner revokes a provider credential
// @Tags providers
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Param credential_id path string true "Credential ID"
// @Success 200 {object} object "revoked successfully"
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /workspaces/{name}/providers/{credential_id} [delete]
func (h *Handler) RevokeProvider(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")
	credentialID, rs := validation.GetUUID(r, "credential_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	_, err := h.commands.RevokeCredential(r.Context(), workspaceName, credentialID)
	if err != nil {
		if errx.IsKind(err, "not_found") {
			fun.NotFound("credential not found").Send(w)
			return
		}
		fun.Error(err).Send(w)
		return
	}

	fun.OK("revoked successfully").Send(w)
}

type SetMarketplaceConfigRequest struct {
	CredentialID uuid.UUID `json:"credential_id" validate:"required"`
	FeeBps       int       `json:"fee_bps" validate:"min=0,max=10000"`
}

// SetMarketplaceConfig godoc
// @Summary Configure marketplace settings for a workspace
// @Description Sets the MP credential and platform fee for marketplace split payments
// @Tags oauth
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Param request body SetMarketplaceConfigRequest true "Marketplace config"
// @Success 200 {object} models.MarketplaceConfig
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 403 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /workspaces/{name}/marketplace [put]
func (h *Handler) SetMarketplaceConfig(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")

	var req SetMarketplaceConfigRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		fun.Error(err).Send(w)
		return
	}

	config, err := h.commands.SetMarketplaceConfig(r.Context(), SetMarketplaceConfigInput{
		WorkspaceName: workspaceName,
		CredentialID:  req.CredentialID,
		FeeBps:        req.FeeBps,
	})
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(config).Send(w)
}

type SetupProviderRequest struct {
	IsMarketplace       bool   `json:"is_marketplace"`
	FeeBps              int    `json:"fee_bps" validate:"min=0,max=10000"`
	ProviderRedirectURL string `json:"provider_redirect_url" validate:"required,url"`
	FinalRedirectURL    string `json:"final_redirect_url" validate:"required,url"`
}

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
// @Param request body SetupProviderRequest true "Setup request"
// @Success 200 {object} BeginOAuthResponse
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /workspaces/{name}/providers/{provider}/setup [post]
func (h *Handler) SetupProvider(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")
	provider := chi.URLParam(r, "provider")

	var req SetupProviderRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		fun.Error(err).Send(w)
		return
	}

	redirectURL, finalRedirectURL, err := h.commands.SetupProvider(r.Context(), SetupProviderInput{
		WorkspaceName:       workspaceName,
		Provider:            provider,
		IsMarketplace:       req.IsMarketplace,
		FeeBps:              req.FeeBps,
		ProviderRedirectURL: req.ProviderRedirectURL,
		FinalRedirectURL:    req.FinalRedirectURL,
	})
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(BeginOAuthResponse{
		RedirectURL:      redirectURL,
		FinalRedirectURL: finalRedirectURL,
	}).Send(w)
}
