package service

import (
	"context"
	"lib/telemetry"
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

	delivery, err := o.Deliveries.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	endpoint, err := o.Endpoints.GetByID(ctx, delivery.EndpointID)
	if err != nil {
		return nil, err
	}
	err = o.Authz.CheckWalletAccess(ctx, ident.Sub.ID, endpoint.WalletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return delivery, nil
}
