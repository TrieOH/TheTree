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
// @Summary Creates a project
// @Description This endpoint creates a project that will consume the Authentication service.
// @Description this project is subjected to limits
// @Tags projects
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 201 {object} dto.ProjectResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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
// @Description
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project to retrieve keys"
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {array} dto.ProjectResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects [get]
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
// @Description
// @Tags projects
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {array} dto.ProjectResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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
// @Description
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project to retrieve keys"
// @Success 200 {object} map[string]any
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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
// @Description
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project to retrieve keys"
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {object} dto.ProjectResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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
// @Description since this is a dangerous action so implement triple confirmation on frontend
// @Description Type the name of the project and hold you are sure button
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project to retrieve keys"
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {string} string "Deleted project"
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
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
