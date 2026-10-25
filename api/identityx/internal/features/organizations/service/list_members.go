package service

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

func (o *Operations) ListMembers(ctx context.Context, orgID uuid.UUID) (members []models.OrganizationMember, err error) {
	err = o.Authz.CheckOrg(ctx, orgID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	members, err = o.Orgs.ListMembers(ctx, orgID)
	if err != nil {
		return nil, err
	}

	return members, nil
}
