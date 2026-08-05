package webhook_endpoints

import (
	"context"
	"lib/telemetry"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (o *Operations) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "DeleteWebhookEndpoint")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	endpoint, err := o.endpoints.GetByID(ctx, id)
	if err != nil {
		return err
	}
	err = o.authz.CheckWalletAccess(ctx, ident.Sub.ID, endpoint.WalletID, models.OrganizationRoleMember)
	if err != nil {
		return err
	}

	return o.endpoints.Delete(ctx, id)
}
