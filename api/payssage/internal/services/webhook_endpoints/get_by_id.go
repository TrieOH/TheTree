package webhook_endpoints

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (o *Operations) GetByID(ctx context.Context, id uuid.UUID) (*models.WebhookEndpoint, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetEndpointByID")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	endpoint, err := o.endpoints.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	err = o.authz.CheckWalletAccess(ctx, ident.Sub.ID, endpoint.WalletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return endpoint, nil
}
