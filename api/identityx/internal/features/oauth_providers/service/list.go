package service

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

func (o *Operations) ListByProject(ctx context.Context, projectID uuid.UUID) ([]models.ProjectOAuthProviders, error) {
	err := o.Authz.CheckProject(ctx, projectID, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Providers.ListByProject(ctx, projectID)
}
