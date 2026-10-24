package service

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (o *Operations) ListByEndpoint(ctx context.Context, endpointID uuid.UUID) ([]models.WebhookDelivery, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListDeliveriesByEndpoint")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	endpoint, err := o.Endpoints.GetByID(ctx, endpointID)
	if err != nil {
		return nil, err
	}
	err = o.Authz.CheckWalletAccess(ctx, ident.Sub.ID, endpoint.WalletID, models.OrganizationRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Deliveries.ListByEndpoint(ctx, endpointID)
}
