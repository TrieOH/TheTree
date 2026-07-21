package queries

import (
	"context"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (q *Queries) ListFromOrg(ctx context.Context, orgID uuid.UUID) ([]models.Wallet, error) {
	ctx, span := q.tracer.Start(ctx, "ListFromOrg")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	org, err := q.orgs.GetByID(ctx, orgID)
	if err != nil {
		return nil, err
	}

	err = q.checkRole(ctx, org, ident.Sub.ID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return q.wallets.ListFromOrg(ctx, org.ID)
}
