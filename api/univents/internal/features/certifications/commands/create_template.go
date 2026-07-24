package commands

import (
	"context"
	"univents/models"
)

func (c *Commands) CreateTemplate(ctx context.Context, input models.CreateCertificationTemplateInput) (*models.CertificationTemplate, error) {
	ctx, span := c.tracer.Start(ctx, "CreateTemplate")
	defer span.End()
	return c.certs.CreateTemplate(ctx, input)
}
