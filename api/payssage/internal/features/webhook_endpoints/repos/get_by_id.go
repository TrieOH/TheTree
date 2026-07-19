package repos

import (
	"context"
	"lib/database"
	"payssage/models"

	"github.com/google/uuid"
)

func (r *repo) GetByID(ctx context.Context, id uuid.UUID) (*models.WebhookEndpoint, error) {
	ctx, span := r.tracer.Start(ctx, "GetByID")
	defer span.End()
	row, err := database.Queries(ctx, r.q).GetWebhookEndpointByID(ctx, id)
	if err != nil {
		return nil, r.dbe(err)
	}
	return new(mapWebhookEndpoint(row)), nil
}
