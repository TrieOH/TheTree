package commands

import (
	"context"

	"github.com/google/uuid"
)

func (c *Commands) SetEditionTemplate(ctx context.Context, editionID uuid.UUID, templateID *uuid.UUID) error {
	ctx, span := c.tracer.Start(ctx, "SetEditionTemplate")
	defer span.End()
	return c.certs.SetEditionTemplate(ctx, editionID, templateID)
}
