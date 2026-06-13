package ports

import (
	"IdentityX/models"
	"context"
)

type ExternalIdentitiesRepo interface {
	GetByProviderAndSubject(ctx context.Context, provider, subject string) (*models.ActorExternalIdentities, error)
	Create(ctx context.Context, identity models.ActorExternalIdentities) (*models.ActorExternalIdentities, error)
	UpdateTokens(ctx context.Context, identity models.ActorExternalIdentities) (*models.ActorExternalIdentities, error)
}
