package repo

import (
	"context"
	"lib/database"
	"payssage/models"

	"github.com/google/uuid"
)

func (r *Repo) IncrementAttempt(ctx context.Context, id uuid.UUID) (*models.WebhookDelivery, error) {
	row, err := database.Queries(ctx, r.q).IncrementDeliveryAttempt(ctx, id)
	if err != nil {
		return nil, r.dbe(err)
	}
	return new(mapWebhookDelivery(row)), nil
}
