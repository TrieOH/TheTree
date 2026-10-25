package service

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

func (o *Operations) SetEnabled(ctx context.Context, id uuid.UUID, enabled bool) (*models.ProjectOAuthProviders, error) {
	row, err := o.Providers.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	err = o.requireProjectAdmin(ctx, row.ProjectID)
	if err != nil {
		return nil, err
	}

	row, err = o.Providers.SetEnabled(ctx, id, enabled)
	if err != nil {
		return nil, err
	}
	return row, nil
}
