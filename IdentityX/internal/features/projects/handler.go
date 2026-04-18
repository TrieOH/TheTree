package projects

import (
	"IdentityX/internal/platform/telemetry"
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/validation"
	"encoding/json"
	"net/http"

	_ "IdentityX/internal/shared/contracts"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/fail/v3"
	"go.uber.org/zap"
)

type Handler struct {
	projects CommandService
}

func NewHandler(
	projects CommandService,
) *Handler {
	return &Handler{projects: projects}
}

type CreateProjectRequest struct {
	ProjectName string          `json:"project_name" validate:"required,max=255"`
	Domain      string          `json:"domain" validate:"required,url"`
	Metadata    json.RawMessage `json:"metadata"`
}

// CreateProject godoc
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
func (handler *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
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
	res, err := handler.projects.Create(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created("Created project").
		WithData(res).
		Send(w)
}

// GetProjectByID godoc
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
func (handler *Handler) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	projectID, rs := validation.GetUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	proj, err := handler.projects.GetByID(ctx, projectID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().
		WithData(proj).
		Send(w)
}

// ListProjects godoc
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
func (handler *Handler) ListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := handler.projects.List(r.Context())
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().
		WithData(projects).
		Send(w)
}

// ListProjectUsers godoc
// @Summary List all users of a project
// @Description Retrieves a list of all users associated with a specific project.
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project"
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {array} contracts.ProjectUser "List of project users"
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

	users, err := handler.projects.ListUsers(r.Context(), projectID)
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
// @Success 200 {object} contracts.ProjectUser "Project user details"
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

	user, err := handler.projects.GetUser(r.Context(), projectID, userID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().
		WithData(user).
		Send(w)
}

// GetProjectJWKS godoc
// @Summary Returns the JWKS for a given project
// @Description Provides the JSON Web Key Set (JWKS) for verifying JWTs issued for a specific project.
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project to retrieve keys"
// @Success 200 {object} object "JSON Web Key Set (JWKS)"
// @Failure 400 {object} contracts.ErrorResponse "Bad Request: Missing project ID"
// @Failure 500 {object} contracts.ErrorResponse "Internal Server Error"
// @Router /projects/{project_id}/.well-known/jwks.json [get]
func (handler *Handler) GetProjectJWKS(w http.ResponseWriter, r *http.Request) {
	projectID, rs := validation.GetUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	jwks, err := handler.projects.GetJWKS(ctx, projectID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	data, err := json.Marshal(jwks)
	if err != nil {
		telemetry.Log().Error("Failed to encode response",
			zap.Error(err),
			zap.String("project_id", projectID.String()),
		)
		apiErr := fail.New(errx.SYSJWKSEncodingFailed).With(err).RecordCtx(ctx)
		resp.FromError(apiErr).Send(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "public, max-age=7200")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(data); err != nil {
		telemetry.Log().Error("Failed to write JWKS response",
			zap.Error(err),
			zap.String("project_id", projectID.String()),
		)
	}
}

type UpdateProjectRequest struct {
	ProjectName string          `json:"project_name" validate:"max=255"`
	Domain      string          `json:"domain" validate:"required,url"`
	Metadata    json.RawMessage `json:"metadata"`
}

// UpdateProjectByID godoc
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
func (handler *Handler) UpdateProjectByID(w http.ResponseWriter, r *http.Request) {
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
	proj, err := handler.projects.Update(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().
		WithData(proj).
		Send(w)
}

// DeleteProjectByID godoc
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
func (handler *Handler) DeleteProjectByID(w http.ResponseWriter, r *http.Request) {
	projectID, rs := validation.GetUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	ctx := r.Context()
	err := handler.projects.Delete(ctx, projectID)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Deleted project").Send(w)
}
