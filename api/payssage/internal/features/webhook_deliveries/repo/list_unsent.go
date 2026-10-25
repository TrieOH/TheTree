package repo

import (
	"context"
	"lib/database"
	"lib/xslices"
	"payssage/internal/sqlc"
	"payssage/models"
	"time"
)

func (r *Repo) ListUnsent(ctx context.Context, olderThan time.Duration, limit int) ([]models.WebhookDelivery, error) {
	rows, err := database.Queries(ctx, r.q).ListUnsentWebhookDeliveries(ctx, sqlc.ListUnsentWebhookDeliveriesParams{
		CreatedBefore: time.Now().Add(-olderThan),
		MaxRows:       limit,
	})
	if err != nil {
		return nil, r.dbe(err)
	}
	return xslices.MapSlice(rows, mapWebhookDelivery), nil
}
