package service

import (
	"context"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) GetByID(ctx context.Context, id uuid.UUID) (*models.Collector, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	collector, err := o.Collectors.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if collector.OrganizationID != nil {
		_, err = o.Authz.Org(ctx, *collector.OrganizationID, models.OrganizationRoleMember)
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
