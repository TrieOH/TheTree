package oauth_providers

import (
	"IdentityX/models"
	"context"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (o *Operations) ListByProject(ctx context.Context, projectID uuid.UUID) ([]models.ProjectOAuthProviders, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListByProject")
	defer span.End()

	err := o.authz.CheckProject(ctx, projectID, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	return o.providers.ListByProject(ctx, projectID)
}
