package ports

import (
	"context"
	"payssage/models"
)

type OAuthStateRepo interface {
	Create(ctx context.Context, state models.OAuthState) (*models.OAuthState, error)
	Get(ctx context.Context, state string) (*models.OAuthState, error)
	Delete(ctx context.Context, state string) error
}
