package service

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

func (o *Operations) List(ctx context.Context, projectID uuid.UUID) ([]models.Actor, error) {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.requireActorReadAccess(ctx, ident, projectID)
	if err != nil {
		return nil, err
	}

	return o.Actors.List(ctx, projectID)
}
