package oauth_handler

import (
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// DeleteMarketplaceConfig godoc
// @Summary Remove marketplace configuration
// @Description Removes the marketplace config for a workspace, reverting to simple mode
// @Tags oauth
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Success 200 {object} object
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /workspaces/{name}/marketplace/{credential_id} [delete]
func (h *Handler) DeleteMarketplaceConfig(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")
	credentialIDStr := chi.URLParam(r, "credential_id")

	credentialID, err := uuid.Parse(credentialIDStr)
	if err != nil {
		resp.BadRequest("invalid credential_id").Send(w)
		return
	}

	if err := h.commands.DeleteMarketplaceConfig(r.Context(), workspaceName, credentialID); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().Send(w)
}
