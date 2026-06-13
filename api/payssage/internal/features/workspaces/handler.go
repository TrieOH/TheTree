package workspaces

import (
	"net/http"

	"payssage/internal/shared/errx"
	"payssage/internal/shared/validation"

	_ "payssage/models"

	"github.com/MintzyG/fun"
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
// @Success 201 {object} models.Workspace "Workspace created successfully"
// @Failure 400 {object} fun.Response
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /workspaces [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateWorkspaceRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		fun.Error(err).Send(w)
		return
	}

	workspace, err := h.commands.Create(r.Context(), req.Name)
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.Created().WithData(workspace).Send(w)
}

// DisableSandbox godoc
// @Summary Disable sandbox mode
// @Description Disables sandbox mode for a workspace, re-enabling real payment processing
// @Tags workspaces
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Success 200 {object} models.Workspace "Sandbox disabled"
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /workspaces/{name}/sandbox/disable [post]
func (h *Handler) DisableSandbox(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	workspace, err := h.commands.DisableSandbox(r.Context(), name)
	if err != nil {
		if errx.IsKind(err, "not_found") {
			fun.NotFound("workspace not found").Send(w)
			return
		}
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(workspace).Send(w)
}

// EnableSandbox godoc
// @Summary Enable sandbox mode
// @Description Enables sandbox mode for a workspace, bypassing real payment processing
// @Tags workspaces
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param name path string true "Workspace name"
// @Success 200 {object} models.Workspace "Sandbox enabled"
// @Failure 401 {object} fun.Response
// @Failure 404 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /workspaces/{name}/sandbox/enable [post]
func (h *Handler) EnableSandbox(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	workspace, err := h.commands.EnableSandbox(r.Context(), name)
	if err != nil {
		if errx.IsKind(err, "not_found") {
			fun.NotFound("workspace not found").Send(w)
			return
		}
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(workspace).Send(w)
}

// List godoc
// @Summary List workspaces
// @Description Lists all workspaces owned by the authenticated user
// @Tags workspaces
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Success 200 {array} models.Workspace "Workspaces retrieved successfully"
// @Failure 401 {object} fun.Response
// @Failure 500 {object} fun.Response
// @Router /workspaces [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	workspaces, err := h.queries.List(r.Context())
	if err != nil {
		fun.Error(err).Send(w)
		return
	}

	fun.OK().WithData(workspaces).Send(w)
}
