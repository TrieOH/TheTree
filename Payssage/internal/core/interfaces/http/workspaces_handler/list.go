package workspaces_handler

import (
	"TriePayments/internal/core/interfaces/http/dto"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

// List godoc
// @Summary List workspaces
// @Description Lists all workspaces owned by the authenticated user
// @Tags workspaces
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Success 200 {array} dto.WorkspaceResponse "Workspaces retrieved successfully"
// @Failure 401 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /workspaces [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	workspaces, err := h.queries.List(r.Context())
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	out := make([]dto.WorkspaceResponse, 0, len(workspaces))
	for _, ws := range workspaces {
		out = append(out, dto.MapWorkspaceResponse(&ws))
	}
}
