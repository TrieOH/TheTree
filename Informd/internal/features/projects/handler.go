package projects

import (
	"TrieForms/internal/shared/validation"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

type Handler struct {
	commands *CommandService
	queries  *QueryService
}

func NewProjectHandler(
	commands *CommandService,
	queries *QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

type CreateProjectRequest struct {
	Name string `json:"name"`
}

// Create godoc
// @Summary Create a project
// @Description Creates a new project for the authenticated user
// @Tags projects
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Param request body CreateProjectRequest true "Project details"
// @Success 201 {object} types.Project "Project created successfully"
// @Failure 400 {object} swag.ErrorResponse
// @Failure 401 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /projects [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.Error(err).Send(w)
		return
	}

	project, err := h.commands.Create(r.Context(), req.Name)
	if err != nil {
		resp.Error(err).Send(w)
		return
	}

	resp.Created().WithData(project).Send(w)
}

// List godoc
// @Summary List projects
// @Description Lists all projects owned by the authenticated user
// @Tags projects
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx"
// @Security Cookie
// @Success 200 {array} types.Project "Projects retrieved successfully"
// @Failure 401 {object} swag.ErrorResponse
// @Failure 500 {object} swag.ErrorResponse
// @Router /projects [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	projects, err := h.queries.List(r.Context())
	if err != nil {
		resp.Error(err).Send(w)
		return
	}

	resp.OK().WithData(projects).Send(w)
}
