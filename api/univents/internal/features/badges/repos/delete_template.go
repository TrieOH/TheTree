package repos

import (
	"context"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (r *Repo) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "BadgeTemplatesRepo.DeleteTemplate")
	defer span.End()
	err := r.q.DeleteBadgeTemplate(ctx, id)
	if err != nil {
		return r.dbe(err)
	}

	return nil
}
