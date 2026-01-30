package handlers

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/adapters/http/validation"
	"GoAuth/internal/ports/inbounds"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

type ScopeHandler struct {
	scopes inbounds.ScopeService
}

func NewScopeHandler(uc inbounds.ScopeService) *ScopeHandler {
	return &ScopeHandler{scopes: uc}
}

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

	in := inbounds.CreateScopeInput{
		ProjectID:  projectID,
		Name:       req.Name,
		ExternalID: req.ExternalID,
	}

	ctx := r.Context()
	scope, err := handler.scopes.Create(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created("Scope Created").WithData(dto.ScopeOutputToScopeResponse(scope)).Send(w)
}

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
