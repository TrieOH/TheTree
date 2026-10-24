package webhook_deliveries

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"
	"payssage/internal/sqlc"
	"payssage/models"
	"time"
)

func (r *Repo) ListUnsent(ctx context.Context, olderThan time.Duration, limit int) ([]models.WebhookDelivery, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListUnsent")
	defer span.End()
	rows, err := database.Queries(ctx, r.q).ListUnsentWebhookDeliveries(ctx, sqlc.ListUnsentWebhookDeliveriesParams{
		CreatedBefore: time.Now().Add(-olderThan),
		MaxRows:       limit,
	})
	if err != nil {
		return nil, r.dbe(err)
	}
	return xslices.MapSlice(rows, mapWebhookDelivery), nil
}
