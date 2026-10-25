package service

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListTemplates(ctx context.Context, editionID uuid.UUID) ([]models.CertificationTemplate, error) {
	return o.Certs.ListTemplates(ctx, editionID)
}
