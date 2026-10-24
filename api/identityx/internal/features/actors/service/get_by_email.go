package service

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

func (o *Operations) GetByEmail(ctx context.Context, email string, projectID uuid.UUID) (*models.Actor, error) {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.requireActorReadAccess(ctx, ident, projectID)
	if err != nil {
		return nil, err
	}

	actor, err := o.Actors.GetByEmail(ctx, email, &projectID)
	if err != nil {
		return nil, err
	}

	return actor, nil
}
