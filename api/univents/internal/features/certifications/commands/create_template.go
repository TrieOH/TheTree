package commands

import (
	"context"
	"univents/contracts"
)

func (c *Commands) CreateTemplate(ctx context.Context, input contracts.CreateCertificationTemplateInput) (*contracts.CertificationTemplate, error) {
	ctx, span := c.tracer.Start(ctx, "CreateTemplate")
	defer span.End()
	return c.certs.CreateTemplate(ctx, input)
}
