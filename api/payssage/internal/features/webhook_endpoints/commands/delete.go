package commands

import (
	"context"
	"lib/telemetry"
	"payssage/internal/authz"
	"payssage/models"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (c *Commands) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "DeleteWebhookEndpoint")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	endpoint, err := c.endpoints.GetByID(ctx, id)
	if err != nil {
		return err
	}
	err = authz.Service.CheckWalletAccess(ctx, ident.Sub.ID, endpoint.WalletID, models.OrganizationRoleMember)
	if err != nil {
		return err
	}

	return c.endpoints.Delete(ctx, id)
}
