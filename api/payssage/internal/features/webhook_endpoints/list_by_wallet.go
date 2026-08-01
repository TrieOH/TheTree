package webhook_endpoints

import (
	"context"
	"lib/telemetry"
	"payssage/internal/authz"
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
	err = authz.Service.CheckWalletAccess(ctx, ident.Sub.ID, walletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.endpoints.ListByWallet(ctx, walletID)
}
