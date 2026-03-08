package workspaces_handler

import (
	"TriePayments/internal/core/interfaces/http/dto"
	"TriePayments/internal/shared/validation"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

// Create godoc
// @Summary Create a workspace
// @Description Creates a new workspace for the authenticated user
// @Tags workspaces
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param request body dto.CreateWorkspaceRequest true "Workspace details"
// @Success 201 {object} dto.WorkspaceResponse "Workspace created successfully"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /workspaces [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateWorkspaceRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	workspace, err := h.commands.Create(r.Context(), req.Name)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created().WithData(dto.MapWorkspaceResponse(workspace)).Send(w)
}
