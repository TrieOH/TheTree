package workspaces_handler

import (
	"TriePayments/internal/core/interfaces/http/dto"
	"TriePayments/internal/shared/errx"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

// DisableSandbox godoc
// @Summary Disable sandbox mode
// @Description Disables sandbox mode for a workspace, re-enabling real payment processing
// @Tags workspaces
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Success 200 {object} dto.WorkspaceResponse "Sandbox disabled"
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /workspaces/{name}/sandbox/disable [post]
func (h *Handler) DisableSandbox(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	workspace, err := h.commands.DisableSandbox(r.Context(), name)
	if err != nil {
		if errx.IsKind(err, "not_found") {
			resp.NotFound("workspace not found").Send(w)
			return
		}
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.MapWorkspaceResponse(workspace)).Send(w)
}
