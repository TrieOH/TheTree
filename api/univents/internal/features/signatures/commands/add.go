package commands

import (
	"context"
	"univents/models"
)

func (c *Commands) Add(ctx context.Context, payload models.AddSignatureInput) (*models.Signature, error) {
	ctx, span := c.tracer.Start(ctx, "Add")
	defer span.End()

	sig := models.Signature{
		EditionID: payload.EditionID,
		Title:     payload.Title,
		URL:       payload.URL,
	}

	return c.signatures.Add(ctx, sig)
}
