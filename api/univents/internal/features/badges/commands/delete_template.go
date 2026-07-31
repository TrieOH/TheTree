package commands

import (
	"context"

	"github.com/google/uuid"
)

func (c *Commands) DeleteTemplate(ctx context.Context, editionID, templateID uuid.UUID) error {
	// Skipping permission checks as requested
	return c.repo.DeleteTemplate(ctx, templateID)
}
