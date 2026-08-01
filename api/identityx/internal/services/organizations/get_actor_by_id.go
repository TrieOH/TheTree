package organizations

import (
	"IdentityX/models"
	"context"
	"lib/telemetry"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) GetActorByID(ctx context.Context, id, orgID, projectID uuid.UUID) (*models.Actor, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetActorByID")
	defer span.End()
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckProject(ctx, ident.Sub.ID, projectID, &orgID, models.ProjectRoleMember)
	if err != nil {
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
