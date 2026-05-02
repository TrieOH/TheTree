package projects

import (
	"IdentityX/internal/shared/contracts"
	"net/http"

	_ "IdentityX/internal/shared/contracts"

	"github.com/MintzyG/fun"
	"github.com/MintzyG/fun/bind"
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

// Create godoc
// @Summary Creates a new project
// @Description Creates a new project that will consume the Authentication service.
// @Tags projects
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param projectInfo body contracts.CreateProjectRequest true "Project creation information"
// @Success 201 {object} contracts.Project "Project created successfully"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Invalid input"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /projects [post]
func (handler *Handler) Create(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	var payload contracts.CreateProjectRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	ctx := r.Context()
	project, err := handler.commands.Create(ctx, payload.ToInput())
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, project)
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
	req := fun.From(r)
	projectID, err := req.Path("project_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	project, err := handler.queries.GetByID(r.Context(), projectID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, project)
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
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, projects)
}

// Update godoc
// @Summary Updates project information
// @Description Updates the name and/or metadata for a specific project.
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project to update"
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param projectInfo body contracts.UpdateProjectRequest true "Project update information"
// @Success 200 {object} contracts.Project "Project updated successfully"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Invalid input or missing project ID"
// @Failure 401 {object} contracts.ErrorResponse "Unauthorized: User not authenticated"
// @Failure 404 {object} contracts.ErrorResponse "Not Found: Project not found"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /projects/{project_id} [patch]
func (handler *Handler) Update(w http.ResponseWriter, r *http.Request) {
	req := fun.From(r)
	projectID, err := req.Path("project_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	var payload contracts.UpdateProjectRequest
	if bind.BailInto(w, req, &payload) {
		return
	}
	project, err := handler.commands.Update(r.Context(), payload.ToInput(projectID))
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, project)
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
	req := fun.From(r)
	projectID, err := req.Path("project_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	err = handler.commands.Delete(r.Context(), projectID)
	if fun.Bail(w, err) {
		return
	}
	fun.OK("Deleted project").Send(w)
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
	req := fun.From(r)
	projectID, err := req.Path("project_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	users, err := handler.queries.ListUsers(r.Context(), projectID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, users)
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
	req := fun.From(r)
	projectID, err := req.Path("project_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	userID, err := req.Path("project_id").UUID()
	if fun.Bail(w, err) {
		return
	}
	user, err := handler.queries.GetUser(r.Context(), projectID, userID)
	if fun.Bail(w, err) {
		return
	}
	fun.Respond(w, user)
}
