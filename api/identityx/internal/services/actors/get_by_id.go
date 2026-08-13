package actors

import (
	"IdentityX/models"
	"context"

	"lib/telemetry"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) GetByID(ctx context.Context, id, projectID uuid.UUID) (*models.Actor, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetByID")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	if err := o.requireActorReadAccess(ctx, ident, projectID); err != nil {
		return nil, err
	}

	actor, err := o.actors.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if actor.ProjectID != nil && *actor.ProjectID != projectID {
		return nil, fun.ErrForbidden("insufficient permissions")
	}

	return actor, nil
}
