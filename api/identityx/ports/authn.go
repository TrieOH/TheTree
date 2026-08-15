package ports

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

type ExternalIdentitiesRepo interface {
	// GetByProviderAndSubject finds the identity for the provider account
	// within the given scope (nil = platform). Identities are scoped by
	// their actor's project, so a project login never sees a platform
	// identity and vice versa.
	GetByProviderAndSubject(ctx context.Context, provider, subject string, projectID *uuid.UUID) (*models.ActorExternalIdentities, error)
	Create(ctx context.Context, identity models.ActorExternalIdentities) (*models.ActorExternalIdentities, error)
	UpdateTokens(ctx context.Context, identity models.ActorExternalIdentities, projectID *uuid.UUID) (*models.ActorExternalIdentities, error)
}
