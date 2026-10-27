package service

import (
	"IdentityX/errs"
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

func (o *Operations) GetActorByID(ctx context.Context, id, projectID uuid.UUID) (*models.Actor, error) {
	err := o.Authz.CheckProject(ctx, projectID, models.ProjectRoleMember)
	if err != nil {
		return nil, err
	}

	actor, err := o.Actors.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if actor.ProjectID != nil && *actor.ProjectID != projectID {
		return nil, errs.ErrActorOutsideProject
	}

	return actor, nil
}
