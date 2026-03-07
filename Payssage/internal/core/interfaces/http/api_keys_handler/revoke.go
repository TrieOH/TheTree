package workspaces_handler

import (
	"TriePayments/internal/shared/validation"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

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
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
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
