package repos

import (
	"context"
	"lib/database"
	"payssage/internal/database/sqlc"
	"payssage/models"
)

func (r *repo) Update(ctx context.Context, params models.UpdateDeliveryParams) (*models.WebhookDelivery, error) {
	ctx, span := r.tracer.Start(ctx, "Update")
	defer span.End()
	row, err := database.Queries(ctx, r.q).UpdateWebhookDelivery(ctx, sqlc.UpdateWebhookDeliveryParams{
		Status:          string(params.Status),
		Attempts:        params.Attempts,
		LastAttemptedAt: params.LastAttemptedAt,
		ResponseStatus:  params.ResponseStatus,
		ResponseBody:    params.ResponseBody,
		ID:              params.ID,
	})
	if err != nil {
		return nil, r.dbe(err)
	}
	return new(mapWebhookDelivery(row)), nil
}
