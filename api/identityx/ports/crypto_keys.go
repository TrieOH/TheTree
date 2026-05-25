package ports

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

type CryptoKeysRepo interface {
	GetActive(ctx context.Context, keyType models.CryptoKeyType, projectID *uuid.UUID) (*models.CryptoKey, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.CryptoKey, error)
}
