package oauth_handler

import (
	"TriePayments/internal/shared/errx"
	"TriePayments/internal/shared/validation"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

// RevokeProvider godoc
// @Summary Revoke a provider credential (owner)
// @Description Workspace owner revokes a provider credential
// @Tags providers
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Param credential_id path string true "Credential ID"
// @Success 200 {object} object "revoked successfully"
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
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
			resp.NotFound("credential not found").Send(w)
			return
		}
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("revoked successfully").Send(w)
}
