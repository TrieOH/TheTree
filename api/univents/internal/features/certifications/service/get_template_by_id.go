package service

import (
	"context"
	"univents/models"
	"uuid"
)

func (o *Operations) GetTemplateByID(ctx context.Context, id uuid.UUID) (*models.CertificationTemplate, error) {
	return o.Certs.GetTemplateByID(ctx, id)
}
