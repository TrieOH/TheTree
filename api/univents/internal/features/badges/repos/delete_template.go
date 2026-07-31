package repos

import (
	"context"

	"github.com/google/uuid"
)

func (r *Repo) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	err := r.q.DeleteBadgeTemplate(ctx, id)
	if err != nil {
		return r.dbe(err)
	}

	return nil
}
