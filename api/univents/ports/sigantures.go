package ports

import (
	"context"
	"univents/contracts"

	"github.com/google/uuid"
)

type SignatureRepo interface {
	Add(ctx context.Context, toAdd contracts.Signature) (*contracts.Signature, error)
	Remove(ctx context.Context, id, editionID uuid.UUID) error
	List(ctx context.Context, editionID uuid.UUID) ([]contracts.Signature, error)
	GetByID(ctx context.Context, id, editionID uuid.UUID) (*contracts.Signature, error)
}
