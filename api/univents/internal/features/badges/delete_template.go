package badges

import (
	"context"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (o *Operations) DeleteTemplate(ctx context.Context, templateID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "BadgeTemplatesCommands.DeleteTemplate")
	defer span.End()
	return o.repo.Delete(ctx, templateID)
}
