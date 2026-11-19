// Package collectors is the collectors feature: its repo, handlers and
// operations. Build wires it.
package collectors

import (
	"context"
	"payssage/errs"
	"payssage/internal/authz"
	"payssage/models"
	"payssage/ports"
	idx "sdk/identityx"
	"uuid"
)

type Deps struct {
	Collectors ports.CollectorRepo
	Orgs       ports.OrganizationRepo
	Authz      *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

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
		return nil, errs.ErrNotResourceOwner
	}

	return collector, nil
}

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

func (o *Operations) ListOwned(ctx context.Context) ([]models.Collector, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	return o.Collectors.ListByOwner(ctx, ident.Sub.ID)
}
