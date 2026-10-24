package service

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (o *Operations) ListByOrg(ctx context.Context, orgID uuid.UUID) ([]models.Intent, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListByOrg")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	org, err := o.Orgs.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	err = o.Authz.CheckOrg(ctx, ident.Sub.ID, org.ID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Intents.ListByOrg(ctx, org.ID)
}
