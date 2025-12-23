package handler

import (
	"GoAuth/internal/models"
	"encoding/json"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/MintzyG/FastUtilitiesNet/validation"
)

// CreateProject godoc
// @Summary Creates a project
// @Description This endpoint creates a project that will consume the Authentication service.
// @Description this project is subjected to limits
// @Tags projects
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 201 {object} models.Project
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /projects [post]
func (h *AuthHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var req models.Project
	if rs := validation.ValidateInto(r, &req); rs != nil {
		rs.Send(w)
		return
	}

	project, rs := h.AuthService.CreateProject(r.Context(), req)
	if rs != nil {
		rs.Send(w)
		return
	}

	resp.Created().WithData(project).Send(w)
}

// GetProjectByID godoc
// @Summary Gets a project by its ID
// @Description
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project to retrieve keys"
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {array} models.Project
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /projects [get]
func (h *AuthHandler) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	projectId := r.PathValue("project_id")
	if projectId == "" {
		resp.BadRequest("missing project id parameter").Send(w)
		return
	}

	project, rs := h.AuthService.GetProjectByID(r.Context(), projectId)
	if rs != nil {
		rs.Send(w)
		return
	}

	resp.OK().WithData(project).Send(w)
}

// ListProjects godoc
// @Summary List all user projects
// @Description
// @Tags projects
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {array} models.Project
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /projects [get]
func (h *AuthHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	projects, rs := h.AuthService.ListProjects(r.Context())
	if rs != nil {
		rs.Send(w)
		return
	}

	resp.OK().WithData(projects).Send(w)
}

// GetProjectKeysByID godoc
// @Summary Returns pub and priv keys for the project
// @Description
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project to retrieve keys"
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Success 200 {array} models.ProjectKeys
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /projects/{project_id}/keys [get]
func (h *AuthHandler) GetProjectKeysByID(w http.ResponseWriter, r *http.Request) {
	projectId := r.PathValue("project_id")
	if projectId == "" {
		resp.BadRequest("missing project id parameter").Send(w)
		return
	}

	keys, rs := h.AuthService.GetProjectKeysByID(r.Context(), r, projectId)
	if rs != nil {
		rs.Send(w)
		return
	}

	resp.OK().WithData(keys).Send(w)
}

// GetProjectJWKS godoc
// @Summary Returns the JWKS for a given project
// @Description
// @Tags projects
// @Accept json
// @Produce json
// @Param project_id path string true "ID of the project to retrieve keys"
// @Success 200 {object} map[string]any
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /projects/{project_id}/.well-known/jwks.json [get]
func (h *AuthHandler) GetProjectJWKS(w http.ResponseWriter, r *http.Request) {
	projectId := r.PathValue("project_id")
	if projectId == "" {
		resp.BadRequest("missing project id parameter").Send(w)
		return
	}

	jwks, rs := h.AuthService.GetProjectJWKS(r.Context(), projectId)
	if rs != nil {
		rs.Send(w)
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
// @Success 200 {object} models.Project
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /projects/{project_id} [patch]
func (h *AuthHandler) UpdateProjectByID(w http.ResponseWriter, r *http.Request) {
	projectId := r.PathValue("project_id")
	if projectId == "" {
		resp.BadRequest("missing project id parameter").Send(w)
		return
	}

	var req models.Project
	if rs := validation.ValidateInto(r, &req); rs != nil {
		rs.Send(w)
		return
	}

	project, rs := h.AuthService.UpdateProjectByID(r.Context(), r, projectId, req)
	if rs != nil {
		rs.Send(w)
		return
	}

	resp.OK().WithData(project).Send(w)
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
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /projects/{project_id} [delete]
func (h *AuthHandler) DeleteProjectByID(w http.ResponseWriter, r *http.Request) {
	projectId := r.PathValue("project_id")
	if projectId == "" {
		resp.BadRequest("missing project id parameter").Send(w)
		return
	}

	rs := h.AuthService.DeleteProjectByID(r.Context(), r, projectId)
	if rs != nil {
		rs.Send(w)
		return
	}

	resp.OK("Deleted project").Send(w)
}
