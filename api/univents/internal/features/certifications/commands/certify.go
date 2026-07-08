package commands

import (
	"context"
	"univents/contracts"
)

func (c *Commands) Certify(ctx context.Context, input contracts.CertifyInput) (*contracts.Certification, error) {
	ctx, span := c.tracer.Start(ctx, "Certify")
	defer span.End()
	return c.certs.Certify(ctx, input)
}
