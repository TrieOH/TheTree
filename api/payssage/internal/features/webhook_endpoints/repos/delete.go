package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "Delete")
	defer span.End()
	err := database.Queries(ctx, r.q).DeleteWebhookEndpoint(ctx, id)
	if err != nil {
		return r.dbe(err)
	}
	return nil
}
