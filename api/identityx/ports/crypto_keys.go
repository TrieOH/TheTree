package ports

import (
	"IdentityX/models"
	"context"
	"lib/crypto"

	"github.com/google/uuid"
)

type CryptoKeysRepo interface {
	GetActive(ctx context.Context, keyType models.CryptoKeyType, projectID *uuid.UUID) (*models.CryptoKey, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.CryptoKey, error)
	GetActiveSigningKeys(ctx context.Context, projectID *uuid.UUID) ([]models.ActiveSigningKey, error)
	Create(ctx context.Context, projectID *uuid.UUID, pair *crypto.KeyPair, keyType string) (*models.CryptoKey, error)
}
