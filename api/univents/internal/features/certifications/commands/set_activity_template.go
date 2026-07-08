package commands

import (
	"context"

	"github.com/google/uuid"
)

func (c *Commands) SetActivityTemplate(ctx context.Context, activityID uuid.UUID, templateID *uuid.UUID) error {
	ctx, span := c.tracer.Start(ctx, "SetActivityTemplate")
	defer span.End()
	return c.certs.SetActivityTemplate(ctx, activityID, templateID)
}
