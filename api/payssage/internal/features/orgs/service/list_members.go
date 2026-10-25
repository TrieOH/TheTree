package service

import (
	"context"
	"payssage/models"

	"github.com/google/uuid"
)

func (o *Operations) ListMembers(ctx context.Context, orgID uuid.UUID) (members []models.OrganizationMember, err error) {
	_, err = o.Authz.Org(ctx, orgID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	members, err = o.Orgs.ListMembers(ctx, orgID)
	if err != nil {
		return nil, err
	}

	return members, nil
}
