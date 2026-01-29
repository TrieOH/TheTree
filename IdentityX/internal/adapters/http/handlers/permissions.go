package handlers

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/adapters/http/validation"
	"GoAuth/internal/ports/inbounds"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
	"github.com/google/uuid"
)

type PermissionHandler struct {
	permission inbounds.PermissionService
}

func NewPermissionHandler(uc inbounds.PermissionService) *PermissionHandler {
	return &PermissionHandler{permission: uc}
}

func (handler *PermissionHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dto.CreatePermissionRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.CreatePermissionInput{
		ProjectID:  &projectID,
		Object:     req.Object,
		Action:     req.Action,
		Conditions: req.Conditions,
	}

	ctx := r.Context()
	perm, err := handler.permission.Create(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created("Permission Created").WithData(dto.PermissionOutputToPermissionResponse(*perm)).Send(w)
}

func (handler *PermissionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	permissionID, rs := getUUID(r, "permission_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	in := inbounds.GetPermissionInput{
		ProjectID:    &projectID,
		PermissionID: permissionID,
	}

	ctx := r.Context()
	perm, err := handler.permission.GetByIDExternal(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.PermissionOutputToPermissionResponse(*perm)).Send(w)
}

func (handler *PermissionHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	object := r.URL.Query().Get("object")
	action := r.URL.Query().Get("action")

	in := inbounds.GetPermissionInput{
		ProjectID: &projectID,
		Object:    &object,
		Action:    &action,
	}

	ctx := r.Context()
	perms, err := handler.permission.ListByProject(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.PermissionOutputSliceToPermissionResponseSlice(perms)).Send(w)
}

func (handler *PermissionHandler) GiveDirect(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	entityID, rs := getUUID(r, "entity_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dto.UserPermissionRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.ManagePermissionInput{
		ProjectID:    &projectID,
		PermissionID: req.PermissionID,
		EntityID:     entityID,
		ScopeID:      req.ScopeID,
	}

	ctx := r.Context()
	err := handler.permission.GiveDirect(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Added permission to user").Send(w)
}

func (handler *PermissionHandler) TakeDirect(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	entityID, rs := getUUID(r, "entity_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dto.UserPermissionRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.ManagePermissionInput{
		ProjectID:    &projectID,
		PermissionID: req.PermissionID,
		EntityID:     entityID,
		ScopeID:      req.ScopeID,
	}

	ctx := r.Context()
	err := handler.permission.TakeDirect(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Removed permission from user").Send(w)
}

func (handler *PermissionHandler) GetEffective(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	entityID, rs := getUUID(r, "entity_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var scopeID *uuid.UUID
	scopeIDStr := r.URL.Query().Get("scope_id")
	if scopeIDStr != "" {
		scopeIDParsed, err := uuid.Parse(scopeIDStr)
		if err != nil {
			resp.FromError(err).Send(w)
			return
		}
		scopeID = &scopeIDParsed
	} else {
		scopeID = nil
	}

	in := inbounds.ManagePermissionInput{
		ProjectID: &projectID,
		ScopeID:   scopeID,
		EntityID:  entityID,
	}

	ctx := r.Context()
	perms, err := handler.permission.GetEffective(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.PermissionOutputSliceToPermissionResponseSlice(perms)).Send(w)
}

func (handler *PermissionHandler) Check(w http.ResponseWriter, r *http.Request) {
	var req dto.CheckRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.CheckPermissionInput{
		ProjectID: req.ProjectID,
		ScopeID:   req.ScopeID,
		EntityID:  req.EntityID,
		Object:    req.Object,
		Action:    req.Action,
	}

	ctx := r.Context()
	hasPermission, err := handler.permission.Check(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	if hasPermission {
		resp.OK("Permission Granted").WithData(map[string]bool{"allowed": hasPermission}).Send(w)
		return
	}

	resp.Forbidden("Permission Denied").WithData(map[string]bool{"allowed": hasPermission}).Send(w)
}
