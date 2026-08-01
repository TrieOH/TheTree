package collectors

import (
	"context"
	"lib/telemetry"
	"payssage/internal/authz"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (o *Operations) ListByOrg(ctx context.Context, orgID uuid.UUID) ([]models.Collector, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListByOrg")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	org, err := o.orgs.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	err = authz.Service.CheckOrg(ctx, ident.Sub.ID, org.ID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.collectors.ListByOrg(ctx, org.ID)
}
