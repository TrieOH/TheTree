package queries

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (q *Queries) GetByID(ctx context.Context, id uuid.UUID) (*models.Collector, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetByID")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	collector, err := q.collectors.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if collector.OrganizationID != nil {
		org, err := q.orgs.GetByID(ctx, *collector.OrganizationID)
		if err != nil {
			return nil, err
		}
		if org.OwnerID == ident.Sub.ID {
			return collector, nil
		}
		_, err = q.orgs.GetMember(ctx, ident.Sub.ID, org.ID)
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
