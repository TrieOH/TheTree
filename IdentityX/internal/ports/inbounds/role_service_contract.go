package inbounds

import (
	"GoAuth/internal/domain/roles"
	"encoding/json"

	"github.com/google/uuid"
)

type RoleInput struct {
	RoleID      uuid.UUID
	ProjectID   *uuid.UUID
	Name        string
	Description *string
	Meta        *json.RawMessage
}

type GetRoleInput struct {
	EntityID  uuid.UUID
	RoleID    uuid.UUID
	ProjectID *uuid.UUID
	Name      string
}

type RolePermissionInput struct {
	ProjectID    *uuid.UUID
	RoleID       uuid.UUID
	PermissionID uuid.UUID
}

type ManageRoleInput struct {
	RoleID    uuid.UUID
	RoleName  string
	EntityID  uuid.UUID
	ScopeID   *uuid.UUID
	ProjectID *uuid.UUID
}

// RoleNotFoundByNameError is returned when a role cannot be found by its name
type RoleNotFoundByNameError struct {
	Name      string
	ProjectID *uuid.UUID
}

func (e *RoleNotFoundByNameError) Error() string {
	if e.ProjectID != nil {
		return "role not found by name: " + e.Name + " in project: " + e.ProjectID.String()
	}
	return "role not found by name: " + e.Name
}

// InvalidRoleNameError is returned when a role name is invalid
type InvalidRoleNameError struct {
	Name string
}

func (e *InvalidRoleNameError) Error() string {
	return "invalid role name: " + e.Name
}

type RoleOutput struct {
	Role roles.Role
}

func RoleToRoleOutput(role roles.Role) *RoleOutput {
	return &RoleOutput{role}
}

func RoleSliceToRoleOutputSlice(roles []roles.Role) []RoleOutput {
	if roles == nil {
		return nil
	}

	out := make([]RoleOutput, 0, len(roles))
	for _, role := range roles {
		out = append(out, RoleOutput{role})
	}
	return out
}
