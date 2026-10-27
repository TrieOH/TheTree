package repo

import (
	"context"
	"lib/database"

	"github.com/google/uuid"
)

func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	err := database.Queries(ctx, r.q).DeleteWebhookEndpoint(ctx, id)
	if err != nil {
		return r.dbe(err)
	}
	return nil
}
