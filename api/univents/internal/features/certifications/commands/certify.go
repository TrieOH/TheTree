package commands

import (
	"context"
	"univents/models"
)

func (c *Commands) Certify(ctx context.Context, input models.CertifyInput) (*models.Certification, error) {
	ctx, span := c.tracer.Start(ctx, "Certify")
	defer span.End()
	return c.certs.Certify(ctx, input)
}
