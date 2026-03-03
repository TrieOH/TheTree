package handlers

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/adapters/http/validation"
	"GoAuth/internal/ports/inbounds"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
)

type RoleHandler struct {
	role inbounds.RoleService
}

func NewRoleHandler(uc inbounds.RoleService) *RoleHandler {
	return &RoleHandler{role: uc}
}

// Create godoc
// @Summary Create a new role
// @Description Creates a new role definition for a project.
// @Tags roles
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param roleInfo body dto.CreateRoleRequest true "Role creation information"
// @Success 201 {object} dto.RoleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/roles [post]
func (handler *RoleHandler) Create(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dto.CreateRoleRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.RoleInput{
		ProjectID:   &projectID,
		Name:        req.Name,
		Description: req.Description,
		Meta:        req.Meta,
	}

	ctx := r.Context()
	role, err := handler.role.Create(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created("Role Created").WithData(dto.RoleOutputToRoleResponse(*role)).Send(w)
}

// GetByID godoc
// @Summary Get role by ID
// @Description Retrieves a role definition by its ID.
// @Tags roles
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param role_id path string true "Role ID"
// @Success 200 {object} dto.RoleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/roles/{role_id} [get]
func (handler *RoleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	roleID, rs := getUUID(r, "role_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	in := inbounds.GetRoleInput{
		ProjectID: &projectID,
		RoleID:    roleID,
	}

	ctx := r.Context()
	role, err := handler.role.GetByIDExternal(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.RoleOutputToRoleResponse(*role)).Send(w)
}

// GetByName godoc
// @Summary Get role by name
// @Description Retrieves a role definition by its name.
// @Tags roles
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param name query string true "Role Name"
// @Success 200 {object} dto.RoleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/roles/search [get]
func (handler *RoleHandler) GetByName(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	name := r.URL.Query().Get("name")

	in := inbounds.GetRoleInput{
		ProjectID: &projectID,
		Name:      name,
	}

	ctx := r.Context()
	role, err := handler.role.GetByName(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.RoleOutputToRoleResponse(*role)).Send(w)
}

// UpdateDescription godoc
// @Summary Update role description
// @Description Updates the description of an existing role.
// @Tags roles
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param role_id path string true "Role ID"
// @Param roleInfo body dto.UpdateRoleDescriptionRequest true "Role update information"
// @Success 200 {object} object "Role updated"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/roles/{role_id}/description [patch]
func (handler *RoleHandler) UpdateDescription(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	roleID, rs := getUUID(r, "role_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dto.UpdateRoleDescriptionRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.RoleInput{
		ProjectID:   &projectID,
		RoleID:      roleID,
		Description: req.Description,
	}

	ctx := r.Context()
	err := handler.role.UpdateDescription(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().Send(w)
}

// UpdateMeta godoc
// @Summary Update role meta
// @Description Updates the meta of an existing role.
// @Tags roles
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param role_id path string true "Role ID"
// @Param roleInfo body dto.UpdateRoleMetaRequest true "Role update information"
// @Success 200 {object} object "Role updated"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/roles/{role_id}/meta [patch]
func (handler *RoleHandler) UpdateMeta(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	roleID, rs := getUUID(r, "role_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	var req dto.UpdateRoleMetaRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.RoleInput{
		ProjectID: &projectID,
		RoleID:    roleID,
		Meta:      req.Meta,
	}

	ctx := r.Context()
	err := handler.role.UpdateMeta(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().Send(w)
}

// Delete godoc
// @Summary Deletes a role
// @Description Deletes a role
// @Tags roles
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param role_id path string true "Role ID"
// @Success 200 {object} object "Role deleted"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/roles/{role_id} [delete]
func (handler *RoleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	roleID, rs := getUUID(r, "role_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	in := inbounds.RoleInput{
		ProjectID: &projectID,
		RoleID:    roleID,
	}

	ctx := r.Context()
	err := handler.role.Delete(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().Send(w)
}

// ListByProject godoc
// @Summary List project roles
// @Description Retrieves all role definitions for a project.
// @Tags roles
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Success 200 {array} dto.RoleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/roles [get]
func (handler *RoleHandler) ListByProject(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	in := inbounds.GetRoleInput{
		ProjectID: &projectID,
	}

	ctx := r.Context()
	role, err := handler.role.ListByProject(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.RoleOutputSliceToRoleResponseSlice(role)).Send(w)
}

// AddPermission godoc
// @Summary Add permission to role
// @Description Associates a permission with a role.
// @Tags roles
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param role_id path string true "Role ID"
// @Param permission_id path string true "Permission ID"
// @Success 200 {object} object "Added permission to role"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/roles/{role_id}/permissions/{permission_id} [post]
func (handler *RoleHandler) AddPermission(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	roleID, rs := getUUID(r, "role_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	permissionID, rs := getUUID(r, "permission_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	in := inbounds.RolePermissionInput{
		ProjectID:    &projectID,
		RoleID:       roleID,
		PermissionID: permissionID,
	}

	ctx := r.Context()
	err := handler.role.AddPermission(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Added permission to role").Send(w)
}

// RemovePermission godoc
// @Summary Remove permission from role
// @Description Removes a permission association from a role.
// @Tags roles
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param role_id path string true "Role ID"
// @Param permission_id path string true "Permission ID"
// @Success 200 {object} object "Removed permission from role"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/roles/{role_id}/permissions/{permission_id} [delete]
func (handler *RoleHandler) RemovePermission(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	roleID, rs := getUUID(r, "role_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	permissionID, rs := getUUID(r, "permission_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	in := inbounds.RolePermissionInput{
		ProjectID:    &projectID,
		RoleID:       roleID,
		PermissionID: permissionID,
	}

	ctx := r.Context()
	err := handler.role.RemovePermission(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Removed permission from role").Send(w)
}

// GetPermissions godoc
// @Summary Get role permissions
// @Description Retrieves all permissions associated with a role.
// @Tags roles
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param role_id path string true "Role ID"
// @Success 200 {array} dto.PermissionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/roles/{role_id}/permissions [get]
func (handler *RoleHandler) GetPermissions(w http.ResponseWriter, r *http.Request) {
	projectID, rs := getUUID(r, "project_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	roleID, rs := getUUID(r, "role_id")
	if rs != nil {
		rs.Send(w)
		return
	}

	in := inbounds.RolePermissionInput{
		ProjectID: &projectID,
		RoleID:    roleID,
	}

	ctx := r.Context()
	permissions, err := handler.role.GetPermissions(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.PermissionOutputSliceToPermissionResponseSlice(permissions)).Send(w)
}

// GiveRole godoc
// @Summary Assign role to user
// @Description Assigns a role to a user (entity) within a specific scope.
// @Tags roles
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param entity_id path string true "Identity ID"
// @Param roleInfo body dto.UserRoleRequest true "Role assignment details"
// @Success 200 {object} object "Added role to user"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/identities/{entity_id}/roles [post]
func (handler *RoleHandler) GiveRole(w http.ResponseWriter, r *http.Request) {
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

	var req dto.UserRoleRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.ManageRoleInput{
		ProjectID: &projectID,
		RoleID:    req.RoleID,
		EntityID:  entityID,
		ScopeID:   req.ScopeID,
	}

	ctx := r.Context()
	err := handler.role.GiveRole(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Added role to user").Send(w)
}

// TakeRole godoc
// @Summary Remove role from user
// @Description Removes a role assignment from a user (entity).
// @Tags roles
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param entity_id path string true "Identity ID"
// @Param roleInfo body dto.UserRoleRequest true "Role revocation details"
// @Success 200 {object} object "Removed role from user"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/identities/{entity_id}/roles [delete]
func (handler *RoleHandler) TakeRole(w http.ResponseWriter, r *http.Request) {
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

	var req dto.UserRoleRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.ManageRoleInput{
		ProjectID: &projectID,
		RoleID:    req.RoleID,
		EntityID:  entityID,
		ScopeID:   req.ScopeID,
	}

	ctx := r.Context()
	err := handler.role.TakeRole(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Removed role from user").Send(w)
}

// GiveRoleByName godoc
// @Summary Assign role to user by name
// @Description Assigns a role to a user (entity) within a specific scope using role name.
// @Tags roles
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param entity_id path string true "Identity ID"
// @Param roleInfo body dto.UserRoleByNameRequest true "Role assignment details"
// @Success 200 {object} object "Added role to user"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/identities/{entity_id}/roles/by-name [post]
func (handler *RoleHandler) GiveRoleByName(w http.ResponseWriter, r *http.Request) {
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

	var req dto.UserRoleByNameRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.ManageRoleInput{
		ProjectID: &projectID,
		RoleName:  req.RoleName,
		EntityID:  entityID,
		ScopeID:   req.ScopeID,
	}

	ctx := r.Context()
	err := handler.role.GiveRoleByName(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Added role to user").Send(w)
}

// TakeRoleByName godoc
// @Summary Remove role from user by name
// @Description Removes a role assignment from a user (entity) using role name.
// @Tags roles
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param entity_id path string true "Identity ID"
// @Param roleInfo body dto.UserRoleByNameRequest true "Role revocation details"
// @Success 200 {object} object "Removed role from user"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/identities/{entity_id}/roles/by-name [delete]
func (handler *RoleHandler) TakeRoleByName(w http.ResponseWriter, r *http.Request) {
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

	var req dto.UserRoleByNameRequest
	if err := validation.ValidateInto(r, &req); err != nil {
		resp.FromError(err).Send(w)
		return
	}

	in := inbounds.ManageRoleInput{
		ProjectID: &projectID,
		RoleName:  req.RoleName,
		EntityID:  entityID,
		ScopeID:   req.ScopeID,
	}

	ctx := r.Context()
	err := handler.role.TakeRoleByName(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK("Removed role from user").Send(w)
}

// GetUserRoles godoc
// @Summary Get user roles
// @Description Retrieves all roles assigned to a user (entity).
// @Tags roles
// @Accept json
// @Produce json
// @Param Cookie header string true "Cookie: access_token=xxx; refresh_token=yyy"
// @Param project_id path string true "Project ID"
// @Param entity_id path string true "Identity ID"
// @Success 200 {array} dto.RoleResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects/{project_id}/identities/{entity_id}/roles [get]
func (handler *RoleHandler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
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

	in := inbounds.GetRoleInput{
		ProjectID: &projectID,
		EntityID:  entityID,
	}

	ctx := r.Context()
	roles, err := handler.role.GetUserRoles(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.RoleOutputSliceToRoleResponseSlice(roles)).Send(w)
}
