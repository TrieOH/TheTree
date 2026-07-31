package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"payssage/models"

	"github.com/google/uuid"
)

func (r *Repo) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.WebhookEndpoint, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListByWallet")
	defer span.End()
	rows, err := database.Queries(ctx, r.q).ListWebhookEndpointsByWallet(ctx, walletID)
	if err != nil {
		return nil, r.dbe(err)
	}
	return xslices.MapSlice(rows, mapWebhookEndpoint), nil
}
