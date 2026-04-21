package api_keys

import (
	"IdentityX/internal/shared/validation"
	"net/http"

	_ "IdentityX/internal/shared/contracts"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

type Handler struct {
	apiKeys CommandService
}

func NewHandler(
	apiKeys CommandService,
) *Handler {
	return &Handler{apiKeys: apiKeys}
}

// RotateApiKey godoc
// @Summary Rotate API key for a project
// @Description Generates a new API key for the project and invalidates the previous one. This is the only time the full key will be shown.
// @Tags api_keys
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project"
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {object} object "API key rotated successfully"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized"
// @Failure 404 {object} contracts.ErrorResponse "Project not found"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /projects/{project_id}/api-keys/rotate [post]
func (handler *Handler) RotateApiKey(w http.ResponseWriter, r *http.Request) {
	projectID, rs := validation.GetUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	apiKey, err := handler.apiKeys.Rotate(r.Context(), projectID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("API key rotated").
		WithData(map[string]interface{}{"api_key": apiKey}).
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
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized"
// @Failure 404 {object} contracts.ErrorResponse "Project not found"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /projects/{project_id}/api-keys [delete]
func (handler *Handler) RevokeApiKey(w http.ResponseWriter, r *http.Request) {
	projectID, rs := validation.GetUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	err := handler.apiKeys.Revoke(r.Context(), projectID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("API key revoked").Send(w)
}
