package ports

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

type SignatureRepo interface {
	Add(ctx context.Context, toAdd models.Signature) (*models.Signature, error)
	Remove(ctx context.Context, id, editionID uuid.UUID) error
	List(ctx context.Context, editionID uuid.UUID) ([]models.Signature, error)
	GetByID(ctx context.Context, id, editionID uuid.UUID) (*models.Signature, error)
}
