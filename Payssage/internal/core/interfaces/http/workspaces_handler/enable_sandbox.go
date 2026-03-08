package workspaces_handler

import (
	"TriePayments/internal/core/interfaces/http/dto"
	"TriePayments/internal/shared/errx"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

// EnableSandbox godoc
// @Summary Enable sandbox mode
// @Description Enables sandbox mode for a workspace, bypassing real payment processing
// @Tags workspaces
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Success 200 {object} dto.WorkspaceResponse "Sandbox enabled"
// @Failure 401 {object} swag.ErrorResponse
// @Failure 404 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /workspaces/{name}/sandbox/enable [post]
func (h *Handler) EnableSandbox(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	workspace, err := h.commands.EnableSandbox(r.Context(), name)
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
