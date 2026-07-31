package commands

import (
	"context"
	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (c *Commands) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := c.tracer.Start(ctx, "DeleteWebhookEndpoint")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	endpoint, err := c.endpoints.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = c.checkWalletAccess(ctx, endpoint.WalletID, ident.Sub.ID)
	if err != nil {
		return err
	}

	return c.endpoints.Delete(ctx, id)
}
