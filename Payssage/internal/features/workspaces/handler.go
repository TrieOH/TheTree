package workspaces

import (
	"net/http"
	"payssage/internal/shared/errx"
	"payssage/internal/shared/validation"

	_ "payssage/internal/shared/contracts"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	commands *CommandService
	queries  *QueryService
}

func NewHandler(
	commands *CommandService,
	queries *QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

type CreateWorkspaceRequest struct {
	Name string `json:"name"`
}

// Create godoc
// @Summary Create a workspace
// @Description Creates a new workspace for the authenticated user
// @Tags workspaces
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param request body CreateWorkspaceRequest true "Workspace details"
// @Success 201 {object} contracts.Workspace "Workspace created successfully"
// @Failure 400 {object} contracts.ErrorResponse
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /workspaces [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateWorkspaceRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	workspace, err := h.commands.Create(r.Context(), req.Name)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created().WithData(workspace).Send(w)
}

// DisableSandbox godoc
// @Summary Disable sandbox mode
// @Description Disables sandbox mode for a workspace, re-enabling real payment processing
// @Tags workspaces
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Success 200 {object} contracts.Workspace "Sandbox disabled"
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
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

	resp.OK().WithData(workspace).Send(w)
}

// EnableSandbox godoc
// @Summary Enable sandbox mode
// @Description Enables sandbox mode for a workspace, bypassing real payment processing
// @Tags workspaces
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Success 200 {object} contracts.Workspace "Sandbox enabled"
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 404 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
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

	resp.OK().WithData(workspace).Send(w)
}

// List godoc
// @Summary List workspaces
// @Description Lists all workspaces owned by the authenticated user
// @Tags workspaces
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Success 200 {array} contracts.Workspace "Workspaces retrieved successfully"
// @Failure 401 {object} contracts.ErrorResponse
// @Failure 500 {object} contracts.ErrorResponse
// @Router /workspaces [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	workspaces, err := h.queries.List(r.Context())
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(workspaces).Send(w)
}
