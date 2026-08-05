package intents

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (o *Operations) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.Intent, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListByWallet")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	err = o.authz.CheckWalletAccess(ctx, ident.Sub.ID, walletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.intents.ListByWallet(ctx, walletID)
}
