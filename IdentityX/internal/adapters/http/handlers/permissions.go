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

// Create godoc
// @Summary Create a new permission
// @Description Creates a new permission definition for a project.
// @Tags permissions
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param permissionInfo body dto.CreatePermissionRequest true "Permission creation information"
// @Success 201 {object} dto.PermissionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/permissions [post]
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
		ProjectID: &projectID,
		Object:    req.Object,
		Action:    req.Action,
		Meta:      req.Meta,
	}

	ctx := r.Context()
	perm, err := handler.permission.Create(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created("Permission Created").WithData(dto.PermissionOutputToPermissionResponse(*perm)).Send(w)
}

// UpdateMeta godoc
// @Summary Update Permission meta
// @Description Updates the meta of an existing permission.
// @Tags permissions
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param permission_id path string true "Permission ID"
// @Param permissionInfo body dto.UpdatePermissionRequest true "Permission update information"
// @Success 200 {object} object "Permission updated"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/permissions/{permission_id}/meta [patch]
func (handler *PermissionHandler) UpdateMeta(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	permID, rs := getUUID(r, "permission_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dto.UpdatePermissionRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.UpdatePermissionInput{
		ProjectID: &projectID,
		ID:        permID,
		Meta:      req.Meta,
	}

	ctx := r.Context()
	err := handler.permission.UpdateMeta(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().Send(w)
}

// Delete godoc
// @Summary delete a Permission
// @Description Deletes a permission
// @Tags permissions
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param permission_id path string true "Permission ID"
// @Success 200 {object} object "Permission deleted"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/permissions/{permission_id} [delete]
func (handler *PermissionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	permID, rs := getUUID(r, "permission_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	in := inbounds.DeletePermissionInput{
		ProjectID: &projectID,
		ID:        permID,
	}

	ctx := r.Context()
	err := handler.permission.Delete(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().Send(w)
}

// GetByID godoc
// @Summary Get permission by ID
// @Description Retrieves a permission definition by its ID.
// @Tags permissions
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param permission_id path string true "Permission ID"
// @Success 200 {object} dto.PermissionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/permissions/{permission_id} [get]
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

// ListByProject godoc
// @Summary List project permissions
// @Description Retrieves all permission definitions for a project, optionally filtered by object and action.
// @Tags permissions
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param object query string false "Filter by object"
// @Param action query string false "Filter by action"
// @Success 200 {array} dto.PermissionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/permissions [get]
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

// GiveDirect godoc
// @Summary Give direct permission to user
// @Description Grants a permission directly to a user (entity) within a specific scope using object:action. Idempotent.
// @Tags permissions
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param entity_id path string true "Identity ID"
// @Param permissionInfo body dto.UserPermissionRequest true "Permission assignment details (object:action)"
// @Success 200 {object} object "Added permission to user"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/identities/{entity_id}/permissions [post]
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
		ProjectID: &projectID,
		Object:    req.Object,
		Action:    req.Action,
		EntityID:  entityID,
		ScopeID:   req.ScopeID,
	}

	ctx := r.Context()
	err := handler.permission.GiveDirect(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Added permission to user").Send(w)
}

// TakeDirect godoc
// @Summary Revoke direct permission from user
// @Description Revokes a directly granted permission from a user (entity) using object:action.
// @Tags permissions
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param entity_id path string true "Identity ID"
// @Param permissionInfo body dto.UserPermissionRequest true "Permission revocation details (object:action)"
// @Success 200 {object} object "Removed permission from user"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/identities/{entity_id}/permissions [delete]
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
		ProjectID: &projectID,
		Object:    req.Object,
		Action:    req.Action,
		EntityID:  entityID,
		ScopeID:   req.ScopeID,
	}

	ctx := r.Context()
	err := handler.permission.TakeDirect(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Removed permission from user").Send(w)
}

// GiveDirectByID godoc
// @Summary Give direct permission to user by ID
// @Description Grants a permission directly to a user (entity) within a specific scope using permission ID. Not idempotent.
// @Tags permissions
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param entity_id path string true "Identity ID"
// @Param permissionInfo body dto.UserPermissionByIDRequest true "Permission assignment details (permission_id)"
// @Success 200 {object} object "Added permission to user"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/identities/{entity_id}/permissions/by-id [post]
func (handler *PermissionHandler) GiveDirectByID(w http.ResponseWriter, r *http.Request) {
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

	var req dto.UserPermissionByIDRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.ManagePermissionByIDInput{
		ProjectID:    &projectID,
		PermissionID: req.PermissionID,
		EntityID:     entityID,
		ScopeID:      req.ScopeID,
	}

	ctx := r.Context()
	err := handler.permission.GiveDirectByID(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Added permission to user").Send(w)
}

// TakeDirectByID godoc
// @Summary Revoke direct permission from user by ID
// @Description Revokes a directly granted permission from a user (entity) using permission ID.
// @Tags permissions
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param entity_id path string true "Identity ID"
// @Param permissionInfo body dto.UserPermissionByIDRequest true "Permission revocation details (permission_id)"
// @Success 200 {object} object "Removed permission from user"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/identities/{entity_id}/permissions/by-id [delete]
func (handler *PermissionHandler) TakeDirectByID(w http.ResponseWriter, r *http.Request) {
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

	var req dto.UserPermissionByIDRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.ManagePermissionByIDInput{
		ProjectID:    &projectID,
		PermissionID: req.PermissionID,
		EntityID:     entityID,
		ScopeID:      req.ScopeID,
	}

	ctx := r.Context()
	err := handler.permission.TakeDirectByID(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Removed permission from user").Send(w)
}

// GetEffective godoc
// @Summary Get effective permissions for user
// @Description Retrieves the list of effective permissions a user has, considering roles and direct assignments.
// @Tags permissions
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param entity_id path string true "Identity ID"
// @Param scope_id query string false "Scope ID to filter effective permissions"
// @Success 200 {array} dto.PermissionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/identities/{entity_id}/permissions [get]
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

// Check godoc
// @Summary Check user permission
// @Description Verifies if a user has a specific permission for an action on a resource.
// @Tags permissions
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param checkInfo body dto.CheckRequest true "Permission check parameters"
// @Success 200 {object} object{allowed=bool} "Permission Granted"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} object{allowed=bool} "Permission Denied"
// @Failure 500 {object} ErrorResponse
// @Router /authz/check [post]
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
		Resource:  req.Resource,
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

// EnsureExists godoc
// @Summary Ensure permissions exist
// @Description Ensures that the specified permissions exist, creating any that are missing. Idempotent operation.
// @Tags permissions
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param permissionInfo body dto.EnsurePermissionsRequest true "Permissions to ensure exist"
// @Success 200 {object} dto.EnsurePermissionsResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/permissions/ensure [post]
func (handler *PermissionHandler) EnsureExists(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dto.EnsurePermissionsRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.EnsurePermissionsInput{
		ProjectID: &projectID,
	}

	in.Permissions = make([]inbounds.PermissionDefinition, 0, len(req.Permissions))
	for _, p := range req.Permissions {
		in.Permissions = append(in.Permissions, inbounds.PermissionDefinition{
			Object: p.Object,
			Action: p.Action,
			Meta:   p.Meta,
		})
	}

	ctx := r.Context()
	results, err := handler.permission.EnsureExists(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	response := dto.EnsurePermissionsResponse{
		Permissions: make([]dto.EnsurePermissionResultDTO, 0, len(results)),
	}
	for _, res := range results {
		response.Permissions = append(response.Permissions, dto.EnsurePermissionResultDTO{
			Object:  res.Object,
			Action:  res.Action,
			Created: res.Created,
		})
	}

	resp.OK().WithData(response).Send(w)
}
