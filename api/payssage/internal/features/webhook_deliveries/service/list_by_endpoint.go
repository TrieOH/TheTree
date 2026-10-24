package service

import (
	"context"
	"lib/telemetry"
	"payssage/models"

	"github.com/google/uuid"
)

func (o *Operations) ListByEndpoint(ctx context.Context, endpointID uuid.UUID) ([]models.WebhookDelivery, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListDeliveriesByEndpoint")
	defer span.End()

	endpoint, err := o.Endpoints.GetByID(ctx, endpointID)
	if err != nil {
		return nil, err
	}
	_, err = o.Authz.Wallet(ctx, endpoint.WalletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Deliveries.ListByEndpoint(ctx, endpointID)
}
