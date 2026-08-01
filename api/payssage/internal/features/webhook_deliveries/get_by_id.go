package webhook_deliveries

import (
	"context"
	"lib/telemetry"
	"payssage/internal/authz"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (o *Operations) GetByID(ctx context.Context, id uuid.UUID) (*models.WebhookDelivery, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetDeliveryByID")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	delivery, err := o.deliveries.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	endpoint, err := o.endpoints.GetByID(ctx, delivery.EndpointID)
	if err != nil {
		return nil, err
	}
	err = authz.Service.CheckWalletAccess(ctx, ident.Sub.ID, endpoint.WalletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return delivery, nil
}
