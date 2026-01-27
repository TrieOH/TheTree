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
	}

	ctx := r.Context()
	role, err := handler.role.Create(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.Created("Role Created").WithData(dto.RoleOutputToRoleResponse(*role)).Send(w)
}

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

	var req dto.UpdateRoleRequest
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
