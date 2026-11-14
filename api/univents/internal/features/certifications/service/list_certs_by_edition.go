package service

import (
	"context"
	"univents/models"
	"uuid"
)

func (o *Operations) ListCertsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Certification, error) {
	return o.Certs.ListByEdition(ctx, editionID)
}
