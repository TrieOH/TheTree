package service

import (
	"context"
	"univents/models"
	"uuid"
)

func (o *Operations) ListTemplates(ctx context.Context, editionID uuid.UUID) ([]models.CertificationTemplate, error) {
	return o.Certs.ListTemplates(ctx, editionID)
}
