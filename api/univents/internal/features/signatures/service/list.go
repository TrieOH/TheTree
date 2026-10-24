package service

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Signature, error) {
	return o.Signatures.ListByEdition(ctx, editionID)
}
