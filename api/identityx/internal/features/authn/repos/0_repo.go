package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type repo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.ExternalIdentitiesRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.ExternalIdentitiesRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("external actor identity"),
	}
}

func mapExternalIdentity(src sqlc.ActorExternalIdentity) models.ActorExternalIdentities {
	return models.ActorExternalIdentities{
		ID:                    src.ID,
		ActorID:               src.ActorID,
		Provider:              models.OAuthProvider(src.Provider),
		Subject:               src.Subject,
		Email:                 src.Email,
		EncryptedAccessToken:  src.EncryptedAccessToken,
		EncryptedRefreshToken: src.EncryptedRefreshToken,
		TokenExpiresAt:        src.TokenExpiresAt,
		CreatedAt:             src.CreatedAt,
		UpdatedAt:             src.UpdatedAt,
	}
}
