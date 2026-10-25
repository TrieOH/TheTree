package service

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

func (o *Operations) ListProjectActors(ctx context.Context, projectID uuid.UUID) ([]models.Actor, error) {
	err := o.Authz.CheckProject(ctx, projectID, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Actors.List(ctx, projectID)
}
