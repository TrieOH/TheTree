package service

import (
	"context"
	"lib/telemetry"
	"payssage/models"

	"github.com/google/uuid"
)

func (o *Operations) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.Intent, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListByWallet")
	defer span.End()

	_, err := o.Authz.Wallet(ctx, walletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Intents.ListByWallet(ctx, walletID)
}
