package inbounds

import (
	"GoAuth/internal/domain/roles"

	"github.com/google/uuid"
)

type RoleInput struct {
	RoleID      uuid.UUID
	ProjectID   *uuid.UUID
	Name        string
	Description *string
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
	EntityID  uuid.UUID
	ScopeID   *uuid.UUID
	ProjectID *uuid.UUID
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

type ErrRoleNotOwned struct {
	Msg string
}

func (e ErrRoleNotOwned) Error() string {
	return e.Msg
}

type ErrProjectUserNotFromProject struct{}

func (e ErrProjectUserNotFromProject) Error() string {
	return "project user not from project"
}
