package queries

import (
	"IdentityX/models"
	"context"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (q *Queries) ListMembers(ctx context.Context, orgID uuid.UUID) (members []models.OrganizationMember, err error) {
	ctx, span := telemetry.StartSpan(ctx, "ListMembers")
	defer span.End()
	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = q.authz.CheckOrg(ctx, ident.Sub.ID, orgID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	members, err = q.orgs.ListMembers(ctx, orgID)
	if err != nil {
		return nil, err
	}

	return members, nil
}
