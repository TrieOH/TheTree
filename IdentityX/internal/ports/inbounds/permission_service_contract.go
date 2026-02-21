package inbounds

import (
	"GoAuth/internal/domain/permissions"
	"encoding/json"

	"github.com/google/uuid"
)

type CreatePermissionInput struct {
	ProjectID *uuid.UUID
	Object    string
	Action    string
	Meta      *json.RawMessage
}

type UpdatePermissionInput struct {
	ID        uuid.UUID
	ProjectID *uuid.UUID
	Meta      *json.RawMessage
}

type DeletePermissionInput struct {
	ID        uuid.UUID
	ProjectID *uuid.UUID
}

type GetPermissionInput struct {
	PermissionID uuid.UUID
	ProjectID    *uuid.UUID
	Object       *string
	Action       *string
}

type ManagePermissionInput struct {
	PermissionID uuid.UUID
	EntityID     uuid.UUID
	ScopeID      *uuid.UUID
	ProjectID    *uuid.UUID
}

type CheckPermissionInput struct {
	EntityID  uuid.UUID
	ProjectID *uuid.UUID
	ScopeID   *uuid.UUID
	Object    string
	Action    string
	Resource  *map[string]interface{}
}

type PermissionOutput struct {
	Permission permissions.Permission
}

func PermissionToPermissionOutput(permission permissions.Permission) *PermissionOutput {
	return &PermissionOutput{permission}
}

func PermissionSliceToPermissionOutputSlice(permissions []permissions.Permission) []PermissionOutput {
	if permissions == nil {
		return nil
	}

	out := make([]PermissionOutput, 0, len(permissions))
	for _, permission := range permissions {
		out = append(out, PermissionOutput{permission})
	}
	return out
}
