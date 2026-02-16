package handlers

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/adapters/http/validation"
	"GoAuth/internal/ports/inbounds"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/google/uuid"
)

type ScopeHandler struct {
	scopes inbounds.ScopeService
}

func NewScopeHandler(uc inbounds.ScopeService) *ScopeHandler {
	return &ScopeHandler{scopes: uc}
}

// Create godoc
// @Summary Create a new scope
// @Description Creates a new scope definition for a project.
// @Tags scopes
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param scopeInfo body dto.CreateScopeRequest true "Scope creation information"
// @Success 201 {object} dto.ScopeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/scopes [post]
func (handler *ScopeHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dto.CreateScopeRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	// Parse optional parent_id
	var parentID *uuid.UUID
	if req.ParentID != nil {
		parsed, err := uuid.Parse(*req.ParentID)
		if err != nil {
			resp.BadRequest("Invalid parent_id format").Send(w)
			return
		}
		parentID = &parsed
	}

	in := inbounds.CreateScopeInput{
		ProjectID:  projectID,
		Name:       req.Name,
		ExternalID: req.ExternalID,
		ParentID:   parentID,
	}

	ctx := r.Context()
	scope, err := handler.scopes.Create(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created("Scope Created").WithData(dto.ScopeOutputToScopeResponse(scope)).Send(w)
}

// GetByID godoc
// @Summary Get scope by ID
// @Description Retrieves a scope definition by its ID.
// @Tags scopes
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param scope_id path string true "Scope ID"
// @Success 200 {object} dto.ScopeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/scopes/{scope_id} [get]
func (handler *ScopeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	scopeID, rs := getUUID(r, "scope_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	in := inbounds.GetScopeInput{
		ProjectID: projectID,
		ScopeID:   scopeID,
	}

	ctx := r.Context()
	scope, err := handler.scopes.GetByIDExternal(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.ScopeOutputToScopeResponse(scope)).Send(w)
}

// GetProjectScopes godoc
// @Summary List project scopes
// @Description Retrieves all scope definitions for a project.
// @Tags scopes
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Success 200 {array} dto.ScopeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/scopes [get]
func (handler *ScopeHandler) GetProjectScopes(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	in := inbounds.GetScopeInput{
		ProjectID: projectID,
	}

	ctx := r.Context()
	scope, err := handler.scopes.GetProjectScopesExternal(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.ScopeOutputSliceToScopeResponseSlice(scope)).Send(w)
}
