package ports

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

type ProjectOAuthProvidersRepo interface {
	Create(ctx context.Context, toCreate models.ProjectOAuthProviders) (*models.ProjectOAuthProviders, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]models.ProjectOAuthProviders, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.ProjectOAuthProviders, error)
	GetByProjectAndProvider(ctx context.Context, projectID uuid.UUID, provider models.OAuthProvider) (*models.ProjectOAuthProviders, error)
	UpdateClientID(ctx context.Context, id uuid.UUID, clientID string) (*models.ProjectOAuthProviders, error)
	UpdateClientSecret(ctx context.Context, id uuid.UUID, encryptedSecret string) (*models.ProjectOAuthProviders, error)
	SetEnabled(ctx context.Context, id uuid.UUID, enabled bool) (*models.ProjectOAuthProviders, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type OAuthLoginStatesRepo interface {
	CreateState(ctx context.Context, state models.OAuthLoginState) (*models.OAuthLoginState, error)
	GetByState(ctx context.Context, state string) (*models.OAuthLoginState, error)
	DeleteState(ctx context.Context, id uuid.UUID) error
}
