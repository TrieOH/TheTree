package repo

import (
	"context"
	"lib/database"
	"lib/xslices"
	"payssage/models"
)

func (r *Repo) ListByProvider(ctx context.Context, provider string) ([]models.WebhookEvent, error) {
	rows, err := database.Queries(ctx, r.q).ListWebhookEventsByProvider(ctx, provider)
	if err != nil {
		return nil, r.dbe(err)
	}
	return xslices.MapSlice(rows, mapWebhookEvent), nil
}
