package handlers

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/ports/inbounds"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

type ApiKeyHandler struct {
	apiKeyService inbounds.ApiKeyService
}

func NewApiKeyHandler(uc inbounds.ApiKeyService) *ApiKeyHandler {
	return &ApiKeyHandler{apiKeyService: uc}
}

// RotateApiKey godoc
// @Summary Rotate API key for a project
// @Description Generates a new API key for the project and invalidates the previous one. This is the only time the full key will be shown.
// @Tags api_keys
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project"
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {object} dto.ApiKeyRotateResponse "API key rotated successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Project not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /projects/{project_id}/api-keys/rotate [post]
func (handler *ApiKeyHandler) RotateApiKey(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	apiKey, err := handler.apiKeyService.Rotate(r.Context(), projectID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("API key rotated").
		WithData(dto.ApiKeyRotateResponse{ApiKey: apiKey}).
		Send(w)
}

// RevokeApiKey godoc
// @Summary Revoke API key for a project
// @Description Deletes the API key for the project, disabling programmatic access.
// @Tags api_keys
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project"
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {object} object "API key revoked successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized"
// @Failure 404 {object} ErrorResponse "Project not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /projects/{project_id}/api-keys [delete]
func (handler *ApiKeyHandler) RevokeApiKey(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	err := handler.apiKeyService.Revoke(r.Context(), projectID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("API key revoked").Send(w)
}
