package repos

import (
	"context"
	"lib/database"
	"lib/xslices"
	"payssage/models"

	"github.com/google/uuid"
)

func (r *repo) ListByWallet(ctx context.Context, walletID uuid.UUID) ([]models.WebhookEvent, error) {
	ctx, span := r.tracer.Start(ctx, "ListByWallet")
	defer span.End()
	rows, err := database.Queries(ctx, r.q).ListWebhookEventsByWallet(ctx, walletID)
	if err != nil {
		return nil, r.dbe(err)
	}
	return xslices.MapSlice(rows, mapWebhookEvent), nil
}
