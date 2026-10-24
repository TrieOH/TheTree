package service

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) GetByID(ctx context.Context, id uuid.UUID) (*models.Signature, error) {
	return o.Signatures.GetByID(ctx, id)
}
