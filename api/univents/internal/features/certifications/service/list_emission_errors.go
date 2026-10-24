package service

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListEmissionErrors(ctx context.Context, editionID uuid.UUID) ([]models.CertEmissionError, error) {
	return o.Certs.ListEmissionErrorsByEdition(ctx, editionID)
}
