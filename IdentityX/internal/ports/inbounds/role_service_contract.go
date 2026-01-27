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
	RoleID    uuid.UUID
	ProjectID *uuid.UUID
	Name      string
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
