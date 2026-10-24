package repo

import (
	"context"
	"lib/database"
	"lib/xslices"
	"payssage/models"

	"github.com/google/uuid"
)

func (r *Repo) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.WebhookEndpoint, error) {
	rows, err := database.Queries(ctx, r.q).ListWebhookEndpointsByWallet(ctx, walletID)
	if err != nil {
		return nil, r.dbe(err)
	}
	return xslices.MapSlice(rows, mapWebhookEndpoint), nil
}
