package repos

import (
	"context"
	"lib/database"
	"payssage/models"

	"github.com/google/uuid"
)

func (r *Repo) MarkDelivered(ctx context.Context, id uuid.UUID) (*models.WebhookDelivery, error) {
	ctx, span := r.tracer.Start(ctx, "MarkDelivered")
	defer span.End()
	row, err := database.Queries(ctx, r.q).MarkDeliveryDelivered(ctx, id)
	if err != nil {
		return nil, r.dbe(err)
	}
	return new(mapWebhookDelivery(row)), nil
}
