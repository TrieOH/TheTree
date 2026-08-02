package collectors

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) GetByID(ctx context.Context, id uuid.UUID) (*models.Collector, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetByID")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	collector, err := o.collectors.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if collector.OrganizationID != nil {
		err = o.authz.CheckOrg(ctx, ident.Sub.ID, *collector.OrganizationID, models.OrganizationRoleMember)
		if err != nil {
			return nil, err
		}
		return collector, nil
	}

	if collector.OwnerID != ident.Sub.ID {
		return nil, fun.ErrForbidden("insufficient permissions")
	}

	return collector, nil
}
