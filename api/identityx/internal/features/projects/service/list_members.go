package service

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

func (o *Operations) ListMembers(ctx context.Context, projectID uuid.UUID) (members []models.ProjectMember, err error) {
	err = o.Authz.CheckProject(ctx, projectID, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	members, err = o.Projects.ListMembers(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return members, nil
}
