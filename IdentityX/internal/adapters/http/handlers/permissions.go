package handlers

import (
	"GoAuth/internal/adapters/http/dto"
	"GoAuth/internal/adapters/http/validation"
	"GoAuth/internal/ports/inbounds"
	"net/http"

	resp "github.com/MintzyG/FastUtilitiesNet/response"
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
	perm, err := handler.permission.ListByProject(ctx, in)
	if err != nil {
		resp.FromError(err).Send(w)
		return
	}

	resp.OK().WithData(dto.PermissionOutputSliceToPermissionResponseSlice(perm)).Send(w)
}
