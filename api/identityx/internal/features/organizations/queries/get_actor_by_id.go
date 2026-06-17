package queries

import (
	"IdentityX/models"
	"context"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (q *Queries) GetActorByID(ctx context.Context, id, orgID, projectID uuid.UUID) (*models.Actor, error) {
	ctx, span := q.tracer.Start(ctx, "OrganizationService.GetByID")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	org, err := q.orgs.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	project, err := q.projects.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project.OrganizationID != nil && *project.OrganizationID != orgID {
		return nil, fun.ErrForbidden("insufficient permissions")
	}

	if ident.Sub.ID != org.OwnerID {
		_, err = q.orgs.GetMember(ctx, ident.Sub.ID, orgID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return nil, err
		}
		if err != nil {
			_, err = q.projects.GetMember(ctx, ident.Sub.ID, projectID)
			if err != nil && !fun.Is(err, fun.CodeNotFound) {
				return nil, err
			}
			if err != nil {
				return nil, fun.ErrForbidden("insufficient permissions")
			}
		}
	}

	actor, err := q.actors.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if actor.ProjectID != nil && *actor.ProjectID != projectID {
		return nil, fun.ErrForbidden("insufficient permissions")
	}

	return actor, nil
}
