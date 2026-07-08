package commands

import (
	"context"
	"univents/contracts"
)

func (c *Commands) Add(ctx context.Context, payload contracts.AddSignatureInput) (*contracts.Signature, error) {
	ctx, span := c.tracer.Start(ctx, "Add")
	defer span.End()

	sig := contracts.Signature{
		EditionID: payload.EditionID,
		Title:     payload.Title,
		URL:       payload.URL,
		PosX:      payload.PosX,
		PosY:      payload.PosY,
	}

	return c.signatures.Add(ctx, sig)
}
