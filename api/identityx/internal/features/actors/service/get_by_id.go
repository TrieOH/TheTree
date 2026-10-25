package service

import (
	"IdentityX/models"
	"context"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) GetByID(ctx context.Context, id, projectID uuid.UUID) (*models.Actor, error) {
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.requireActorReadAccess(ctx, ident, projectID)
	if err != nil {
		return nil, err
	}

	actor, err := o.Actors.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if actor.ProjectID != nil && *actor.ProjectID != projectID {
		return nil, fun.ErrForbidden("insufficient permissions")
	}

	return actor, nil
}
