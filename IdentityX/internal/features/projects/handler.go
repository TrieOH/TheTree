package projects

import (
	"IdentityX/internal/shared/validation"
	"encoding/json"
	"net/http"

	_ "IdentityX/internal/shared/contracts"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

type Handler struct {
	commands CommandService
	queries  QueryService
}

func NewHandler(
	commands CommandService,
	queries QueryService,
) *Handler {
	return &Handler{
		commands: commands,
		queries:  queries,
	}
}

type CreateProjectRequest struct {
	ProjectName string          `json:"project_name" validate:"required,max=255"`
	Domain      string          `json:"domain" validate:"required,url"`
	Metadata    json.RawMessage `json:"metadata"`
}

// Create godoc
// @Summary Creates a new project
// @Description Creates a new project that will consume the Authentication service.
// @Tags projects
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param projectInfo body CreateProjectRequest true "Project creation information"
// @Success 201 {object} contracts.Project "Project created successfully"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Invalid input"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /projects [post]
func (handler *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateProjectRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := ProjectServiceInput{
		ProjectName: req.ProjectName,
		Metadata:    req.Metadata,
		Domain:      req.Domain,
	}

	ctx := r.Context()
	res, err := handler.commands.Create(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created("Created project").
		WithData(res).
		Send(w)
}

// GetByID godoc
// @Summary Gets a project by its ID
// @Description Retrieves details of a specific project by its ID.
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project to retrieve"
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {object} contracts.Project "Project details"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Missing project ID"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: User not authenticated"
// @Failure 404 {object} contracts.ErrorResponse "Not Found: Project not found"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /projects/{project_id} [get]
func (handler *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	projectID, rs := validation.GetUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	proj, err := handler.queries.GetByID(ctx, projectID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().
		WithData(proj).
		Send(w)
}

// List godoc
// @Summary List all user projects
// @Description Retrieves a list of all projects associated with the authenticated user.
// @Tags projects
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {array} contracts.Project "List of user projects"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /projects [get]
func (handler *Handler) List(w http.ResponseWriter, r *http.Request) {
	projects, err := handler.queries.List(r.Context())
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().
		WithData(projects).
		Send(w)
}

type UpdateProjectRequest struct {
	ProjectName string          `json:"project_name" validate:"max=255"`
	Domain      string          `json:"domain" validate:"required,url"`
	Metadata    json.RawMessage `json:"metadata"`
}

// Update godoc
// @Summary Updates project information
// @Description Updates the name and/or metadata for a specific project.
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project to update"
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param projectInfo body UpdateProjectRequest true "Project update information"
// @Success 200 {object} contracts.Project "Project updated successfully"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Invalid input or missing project ID"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: User not authenticated"
// @Failure 404 {object} contracts.ErrorResponse "Not Found: Project not found"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /projects/{project_id} [patch]
func (handler *Handler) Update(w http.ResponseWriter, r *http.Request) {
	projectID, rs := validation.GetUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req UpdateProjectRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := ProjectServiceInput{
		ProjectID:   projectID,
		ProjectName: req.ProjectName,
		Domain:      req.Domain,
		Metadata:    req.Metadata,
	}

	ctx := r.Context()
	proj, err := handler.commands.Update(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().
		WithData(proj).
		Send(w)
}

// Delete godoc
// @Summary Deletes a user project
// @Description Deletes a specific project by its ID. This is a dangerous action and requires careful confirmation.
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project to delete"
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {object} object "Project deleted successfully"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Missing project ID"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: User not authenticated"
// @Failure 404 {object} contracts.ErrorResponse "Not Found: Project not found"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /projects/{project_id} [delete]
func (handler *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	projectID, rs := validation.GetUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	err := handler.commands.Delete(ctx, projectID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Deleted project").Send(w)
}

// ListProjectUsers godoc
// @Summary List all users of a project
// @Description Retrieves a list of all users associated with a specific project.
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project"
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {array} contracts.User "List of project users"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Missing project ID"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: User not authenticated"
// @Failure 404 {object} contracts.ErrorResponse "Not Found: Project not found"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /projects/{project_id}/users [get]
func (handler *Handler) ListProjectUsers(w http.ResponseWriter, r *http.Request) {
	projectID, rs := validation.GetUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	users, err := handler.queries.ListUsers(r.Context(), projectID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().
		WithData(users).
		Send(w)
}

// GetProjectUserByID godoc
// @Summary Gets a project user by its ID
// @Description Retrieves details of a specific user associated with a specific project.
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project"
// @Param user_id path string true "ID of the user"
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {object} contracts.User "Project user details"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Missing project or user ID"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: User not authenticated"
// @Failure 404 {object} contracts.ErrorResponse "Not Found: User or project not found"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /projects/{project_id}/users/{user_id} [get]
func (handler *Handler) GetProjectUserByID(w http.ResponseWriter, r *http.Request) {
	projectID, rs := validation.GetUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	userID, rs := validation.GetUUID(r, "user_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	user, err := handler.queries.GetUser(r.Context(), projectID, userID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().
		WithData(user).
		Send(w)
}
