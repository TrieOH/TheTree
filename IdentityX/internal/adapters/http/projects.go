package http

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/application/project"
	"encoding/json"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/FastUtilitiesNet/validation"
	"github.com/go-chi/chi/v5"
)

type ProjectHandler struct {
	uc *project.UseCase
}

func NewProjectHandler(uc *project.UseCase) *ProjectHandler {
	return &ProjectHandler{uc: uc}
}

// CreateProject godoc
// @Summary Creates a new project
// @Description Creates a new project that will consume the Authentication service.
// @Tags projects
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param projectInfo body dto.CreateProjectRequest true "Project creation information"
// @Success 201 {object} dto.ProjectResponse "Project created successfully"
// @Failure 400 {object} ErrorResponse "Bad Request: Invalid input"
// @Failure 401 {object} ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /projects [post]
func (ph *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProjectRequest
	if rs := validation.ValidateInto(r, &req); rs != nil {
		rs.Send(w)
		return
	}

	in := project.CreateProjectInput{
		ProjectName: req.ProjectName,
		Metadata:    req.Metadata,
	}

	ctx := r.Context()
	res, err := ph.uc.CreateProject(ctx, in)
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	resp.Created("Created project").
		WithData(dto.ProjectToResponse(res)).
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
// @Success 200 {object} dto.ProjectResponse "Project details"
// @Failure 400 {object} ErrorResponse "Bad Request: Missing project ID"
// @Failure 401 {object} ErrorResponse "Unauthorized: User not authenticated"
// @Failure 404 {object} ErrorResponse "Not Found: Project not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /projects/{project_id} [get]
func (ph *ProjectHandler) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "project_id")
	if projectID == "" {
		resp.BadRequest("missing project id parameter").Send(w)
		return
	}

	ctx := r.Context()
	proj, err := ph.uc.GetProjectByID(ctx, projectID)
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	resp.OK().
		WithData(dto.ProjectToResponse(proj)).
		Send(w)
}

// ListProjects godoc
// @Summary List all user projects
// @Description Retrieves a list of all projects associated with the authenticated user.
// @Tags projects
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {array} dto.ProjectResponse "List of user projects"
// @Failure 401 {object} ErrorResponse "Unauthorized: User not authenticated"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /projects [get]
func (ph *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := ph.uc.ListProjects(r.Context())
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	resp.OK().
		WithData(dto.ProjectSliceToProjectResponseSlice(projects)).
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
// @Failure 400 {object} ErrorResponse "Bad Request: Missing project ID"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /projects/{project_id}/.well-known/jwks.json [get]
func (ph *ProjectHandler) GetProjectJWKS(w http.ResponseWriter, r *http.Request) {
	projectId := chi.URLParam(r, "project_id")
	if projectId == "" {
		resp.BadRequest("missing project id parameter").Send(w)
		return
	}

	jwks, err := ph.uc.GetProjectJWKS(r.Context(), projectId)
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"keys": []any{jwks},
	})
}

// UpdateProjectByID godoc
// @Summary Updates project information
// @Description Updates the name and/or metadata for a specific project.
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project to update"
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param projectInfo body dto.UpdateProjectRequest true "Project update information"
// @Success 200 {object} dto.ProjectResponse "Project updated successfully"
// @Failure 400 {object} ErrorResponse "Bad Request: Invalid input or missing project ID"
// @Failure 401 {object} ErrorResponse "Unauthorized: User not authenticated"
// @Failure 404 {object} ErrorResponse "Not Found: Project not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /projects/{project_id} [patch]
func (ph *ProjectHandler) UpdateProjectByID(w http.ResponseWriter, r *http.Request) {
	projectID := chi.URLParam(r, "project_id")
	if projectID == "" {
		resp.BadRequest("missing project id parameter").Send(w)
		return
	}

	var req dto.UpdateProjectRequest
	if rs := validation.ValidateInto(r, &req); rs != nil {
		rs.Send(w)
		return
	}

	in := project.UpdateProjectInput{
		ProjectID:   projectID,
		ProjectName: req.ProjectName,
		Metadata:    req.Metadata,
	}

	ctx := r.Context()
	proj, err := ph.uc.UpdateProjectByID(ctx, in)
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	resp.OK().
		WithData(dto.ProjectToResponse(proj)).
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
// @Failure 400 {object} ErrorResponse "Bad Request: Missing project ID"
// @Failure 401 {object} ErrorResponse "Unauthorized: User not authenticated"
// @Failure 404 {object} ErrorResponse "Not Found: Project not found"
// @Failure 500 {object} ErrorResponse "Internal Server Error"
// @Router /projects/{project_id} [delete]
func (ph *ProjectHandler) DeleteProjectByID(w http.ResponseWriter, r *http.Request) {
	projectId := chi.URLParam(r, "project_id")
	if projectId == "" {
		resp.BadRequest("missing project id parameter").Send(w)
		return
	}

	ctx := r.Context()
	err := ph.uc.DeleteProjectByID(ctx, projectId)
	if err != nil {
		ErrToResp(err).Send(w)
		return
	}

	resp.OK("Deleted project").Send(w)
}
