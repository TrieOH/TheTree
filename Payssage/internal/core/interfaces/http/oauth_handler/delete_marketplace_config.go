package oauth_handler

import (
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
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
// @Router /workspaces/{name}/marketplace [delete]
func (h *Handler) DeleteMarketplaceConfig(w http.ResponseWriter, r *http.Request) {
	workspaceName := chi.URLParam(r, "name")

	if err := h.commands.DeleteMarketplaceConfig(r.Context(), workspaceName); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().Send(w)
}
