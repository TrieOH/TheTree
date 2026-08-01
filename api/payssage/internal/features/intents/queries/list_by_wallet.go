package queries

import (
	"context"
	"lib/telemetry"
	"payssage/internal/authz"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (q *Queries) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.Intent, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListByWallet")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	err = authz.Service.CheckWalletAccess(ctx, ident.Sub.ID, walletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return q.intents.ListByWallet(ctx, walletID)
}
