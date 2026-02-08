package inbounds

import (
	"GoAuth/internal/domain/authz"
	"context"

	"github.com/google/uuid"
)

type ApiKeyService interface {
	Rotate(ctx context.Context, projectID uuid.UUID) (string, error)
	Revoke(ctx context.Context, projectID uuid.UUID) error
	Authenticate(ctx context.Context, apiKey string) (*authz.Principal, error)
}
