package repos

import (
	sqlc2 "Informd/internal/sqlc"
	"Informd/models"
	"Informd/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
)

type repo struct {
	q      *sqlc2.Queries
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.NamespaceRepo = (*repo)(nil)

func NewRepo(q *sqlc2.Queries, tracer trace.Tracer) ports.NamespaceRepo {
	return &repo{
		q:      q,
		tracer: tracer,
		dbe:    database.NewErrorHandler("namespace"),
	}
}

func mapNamespace(src sqlc2.Namespace) models.Namespace {
	return models.Namespace{
		ID:        src.ID,
		OwnerID:   src.OwnerID,
		Name:      src.Name,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
	}
}

func mapNamespaceMember(src sqlc2.NamespaceMember) models.NamespaceMember {
	return models.NamespaceMember{
		UserID:      src.UserID,
		NamespaceID: src.NamespaceID,
		Role:        models.NamespaceMemberRole(src.Role),
		AddedAt:     src.AddedAt,
		AddedBy:     src.AddedBy,
	}
}
