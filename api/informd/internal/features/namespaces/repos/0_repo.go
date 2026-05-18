package repos

import (
	"Informd/internal/database/sqlc"
	"Informd/models"
	"Informd/ports"
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

var _ ports.NamespaceRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.NamespaceRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("namespace"),
	}
}

func mapNamespace(src sqlc.Namespace) models.Namespace {
	return models.Namespace{
		ID:        src.ID,
		OwnerID:   src.OwnerID,
		Name:      src.Name,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
	}
}

func mapNamespaceMember(src sqlc.NamespaceMember) models.NamespaceMember {
	return models.NamespaceMember{
		UserID:      src.UserID,
		NamespaceID: src.NamespaceID,
		Role:        models.NamespaceMemberRole(src.Role),
		AddedAt:     src.AddedAt,
		AddedBy:     src.AddedBy,
	}
}
