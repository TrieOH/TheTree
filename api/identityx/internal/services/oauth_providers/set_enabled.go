package oauth_providers

import (
	"IdentityX/models"
	"context"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (o *Operations) SetEnabled(ctx context.Context, id uuid.UUID, enabled bool) (*models.ProjectOAuthProviders, error) {
	ctx, span := telemetry.StartSpan(ctx, "SetEnabled")
	defer span.End()

	row, err := o.providers.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	err = o.requireProjectAdmin(ctx, row.ProjectID)
	if err != nil {
		return nil, err
	}

	row, err = o.providers.SetEnabled(ctx, id, enabled)
	if err != nil {
		return nil, err
	}
	return row, nil
}
