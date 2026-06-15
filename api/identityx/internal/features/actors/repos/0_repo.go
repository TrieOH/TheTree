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

var _ ports.ActorRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.ActorRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("actor"),
	}
}

func mapActor(src sqlc.Actor) models.Actor {
	return models.Actor{
		ID:           src.ID,
		ProjectID:    src.ProjectID,
		AuthMethod:   models.AuthMethod(src.AuthMethod),
		VerifiedAt:   src.VerifiedAt,
		PasswordHash: src.PasswordHash,
		Email:        src.Email,
		Type:         models.ActorType(src.Type),
		Metadata:     src.Metadata,
		CreatedAt:    src.CreatedAt,
		UpdatedAt:    src.UpdatedAt,
		DeletedAt:    src.DeletedAt,
	}
}
