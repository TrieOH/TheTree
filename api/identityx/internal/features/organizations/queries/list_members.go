package queries

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

func (q *Queries) ListMembers(ctx context.Context, orgID uuid.UUID) (members []models.OrganizationMember, err error) {
	ctx, span := q.tracer.Start(ctx, "OrganizationService.GetMembers")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	var org *models.Organization
	org, err = q.orgs.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	if ident.Sub.ID != org.OwnerID {
		_, err = q.orgs.GetMember(ctx, ident.Sub.ID, orgID)
		if err != nil {
			return nil, err
		}
	}

	members, err = q.orgs.ListMembers(ctx, orgID)
	if err != nil {
		return nil, err
	}

	return members, nil
}
