package ports

import (
	"context"
	"univents/models"

	"github.com/google/uuid"
)

type SignatureRepo interface {
	Create(ctx context.Context, toCreate *models.Signature) (*models.Signature, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.Signature, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Signature, error)
}

type SignatureRequestRepo interface {
	CreateRequest(ctx context.Context, toCreate *models.SignatureRequest) (*models.SignatureRequest, error)
	GetRequestByIdempotencyKey(ctx context.Context, idempotencyKey uuid.UUID) (*models.SignatureRequest, error)
	GetRequestByID(ctx context.Context, id uuid.UUID) (*models.SignatureRequest, error)
	ListRequestsByEdition(ctx context.Context, editionID uuid.UUID) ([]models.SignatureRequest, error)
	CompleteRequest(ctx context.Context, id, signatureID uuid.UUID) error
	CancelRequest(ctx context.Context, id uuid.UUID, reason *string) error
}
