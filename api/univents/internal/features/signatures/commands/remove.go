package commands

import (
	"context"
	"lib/objectstorage"

	"github.com/google/uuid"
)

func (c *Commands) Remove(ctx context.Context, id, editionID uuid.UUID) error {
	ctx, span := c.tracer.Start(ctx, "Remove")
	defer span.End()

	sig, err := c.signatures.GetByID(ctx, id, editionID)
	if err != nil {
		return err
	}

	bucket, key, err := objectstorage.ParseURL(sig.URL)
	if err != nil {
		return err
	}

	err = c.obj.RemoveObject(ctx, bucket, key)
	if err != nil {
		return err
	}

	return c.signatures.Remove(ctx, id, editionID)
}
