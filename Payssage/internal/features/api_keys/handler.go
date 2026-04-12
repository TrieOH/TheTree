package api_keys

import (
	"net/http"
	"payssage/internal/shared/contracts"
	"payssage/internal/shared/validation"

	_ "payssage/internal/shared/contracts"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
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

type CreateAPIKeyRequest struct {
	Name string `json:"name"`
}

type CreateAPIKeyResponse struct {
	ApiKey *contracts.APIKey
	Key    string `json:"key"` // only returned once
}

// Create godoc
// @Summary Create an API key
// @Description Creates a new API key for the given workspace. The raw key is returned once and never stored.
// @Tags api_keys
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Param request body CreateAPIKeyRequest true "API key details"
// @Success 201 {object} CreateAPIKeyResponse "API key created successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /workspaces/{name}/keys [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")

	var req CreateAPIKeyRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	rawKey, apiKey, err := h.commands.Create(r.Context(), workspaceName, req.Name)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created().WithData(CreateAPIKeyResponse{
		ApiKey: apiKey,
		Key:    rawKey,
	}).Send(w)
}

// ListAPIKeys godoc
// @Summary List API keys
// @Description Lists all API keys for the given workspace (raw keys are never returned)
// @Tags api_keys
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Success 200 {array} contracts.APIKey "API keys retrieved successfully"
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /workspaces/{name}/keys [get]
func (h *Handler) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")

	keys, err := h.queries.List(r.Context(), workspaceName)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(keys).Send(w)
}

// RevokeAPIKey godoc
// @Summary Revoke an API key
// @Description Revokes the given API key, immediately invalidating it
// @Tags api_keys
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Param id path string true "API key ID"
// @Success 200 {object} object "Key revoked"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /workspaces/{name}/keys/{id} [delete]
func (h *Handler) RevokeAPIKey(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")

	keyID, rs := validation.GetUUID(r, "id")
	if rs != nil {
		rs.Send(w)
		return
	}

	if err := h.commands.RevokeAPIKey(r.Context(), workspaceName, keyID); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("key revoked").Send(w)
}
