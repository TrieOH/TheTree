package service

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (o *Operations) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.WebhookEndpoint, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListEndpointsByWallet")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	err = o.Authz.CheckWalletAccess(ctx, ident.Sub.ID, walletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Endpoints.ListByWallet(ctx, walletID)
}
