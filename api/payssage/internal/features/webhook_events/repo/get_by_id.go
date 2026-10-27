package repo

import (
	"context"
	"lib/database"
	"payssage/models"

	"github.com/google/uuid"
)

func (r *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.WebhookEvent, error) {
	row, err := database.Queries(ctx, r.q).GetWebhookEventByID(ctx, id)
	if err != nil {
		return nil, r.dbe(err)
	}
	return new(mapWebhookEvent(row)), nil
}
