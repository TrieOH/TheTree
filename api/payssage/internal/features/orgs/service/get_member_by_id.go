package service

import (
	"context"
	"payssage/models"

	"github.com/google/uuid"
)

func (o *Operations) GetMemberByID(ctx context.Context, id, orgID uuid.UUID) (*models.OrganizationMember, error) {
	_, err := o.Authz.Org(ctx, orgID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	actor, err := o.IDX.Actors.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	member, err := o.Orgs.GetMember(ctx, actor.ID, orgID)
	if err != nil {
		return nil, err
	}

	return member, nil
}
