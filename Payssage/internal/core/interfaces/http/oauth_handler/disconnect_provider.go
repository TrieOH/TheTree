package oauth_handler

import (
	"TriePayments/internal/shared/errx"
	"TriePayments/internal/shared/validation"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

// DisconnectProvider godoc
// @Summary Disconnect a provider credential (seller)
// @Description Called via API key by Univents when a seller clicks Disconnect
// @Tags providers
// @Param X-API-Key header string true "X-API-Key: tp_xxxxxxxx"
// @Security APIKey
// @Param name path string true "Workspace name"
// @Param credential_id path string true "Credential ID"
// @Success 200 {object} object "disconnected successfully"
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
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
			resp.NotFound("credential not found").Send(w)
			return
		}
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("disconnected successfully").Send(w)
}
