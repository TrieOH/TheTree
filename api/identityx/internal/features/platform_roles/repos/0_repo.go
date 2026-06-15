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

var _ ports.PlatformRolesRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.PlatformRolesRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("platform role"),
	}
}

func mapPlatformRole(src sqlc.PlatformRole) models.PlatformRoleRelation {
	return models.PlatformRoleRelation{
		ActorID:   src.ActorID,
		Role:      models.PlatformRole(src.Role),
		Metadata:  src.Metadata,
		CreatedAt: src.CreatedAt,
	}
}
