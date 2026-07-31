package commands

import (
	"context"
	"lib/telemetry"

	"github.com/google/uuid"
)

func (c *Commands) DeleteTemplate(ctx context.Context, editionID, templateID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "BadgeTemplatesCommands.DeleteTemplate")
	defer span.End()
	return c.repo.DeleteTemplate(ctx, templateID)
}
