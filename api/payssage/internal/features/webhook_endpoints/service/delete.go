package service

import (
	"context"
	"lib/telemetry"
	"payssage/models"

	"github.com/google/uuid"
)

func (o *Operations) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "DeleteWebhookEndpoint")
	defer span.End()

	endpoint, err := o.Endpoints.GetByID(ctx, id)
	if err != nil {
		return err
	}
	_, err = o.Authz.Wallet(ctx, endpoint.WalletID, models.OrganizationRoleMember)
	if err != nil {
		return err
	}

	return o.Endpoints.Delete(ctx, id)
}
