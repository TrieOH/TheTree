package repos

import (
	"context"
	"lib/database"
	"lib/xslices"
	"payssage/models"
)

func (r *repo) ListByProvider(ctx context.Context, provider string) ([]models.WebhookEvent, error) {
	ctx, span := r.tracer.Start(ctx, "ListByProvider")
	defer span.End()
	rows, err := database.Queries(ctx, r.q).ListWebhookEventsByProvider(ctx, provider)
	if err != nil {
		return nil, r.dbe(err)
	}
	return xslices.MapSlice(rows, mapWebhookEvent), nil
}
