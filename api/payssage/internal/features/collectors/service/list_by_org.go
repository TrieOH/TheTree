package service

import (
	"context"
	"payssage/models"

	"github.com/google/uuid"
)

func (o *Operations) ListByOrg(ctx context.Context, orgID uuid.UUID) ([]models.Collector, error) {
	org, err := o.Orgs.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Org(ctx, org.ID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Collectors.ListByOrg(ctx, org.ID)
}
