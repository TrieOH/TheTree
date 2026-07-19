package repos

import (
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (r *repo) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := r.tracer.Start(ctx, "Delete")
	defer span.End()
	err := database.Queries(ctx, r.q).DeleteWebhookEndpoint(ctx, id)
	if err != nil {
		return r.dbe(err)
	}
	return nil
}
