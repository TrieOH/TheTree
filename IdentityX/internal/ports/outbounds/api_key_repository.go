package outbounds

import (
	"GoAuth/internal/domain/apikey"
	"context"

	"github.com/google/uuid"
)

type ApiKeyRepository interface {
	Upsert(ctx context.Context, key apikey.ApiKey) error
	GetByProjectID(ctx context.Context, projectID uuid.UUID) (*apikey.ApiKey, error)
	Delete(ctx context.Context, projectID uuid.UUID) error
}
