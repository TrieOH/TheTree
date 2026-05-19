package forms

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

var _ ports.FormsRepo = (*repo)(nil)

func NewFormRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.FormsRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("form"),
	}
}

func mapForm(src sqlc.Form) models.Form {
	return models.Form{
		ID:          src.ID,
		NamespaceID: src.NamespaceID,
		OwnerID:     src.OwnerID,
		Title:       src.Name,
		Status:      models.FormStatus(src.Status),
		OpenedAt:    src.OpenedAt,
		ClosedAt:    src.ClosedAt,
		ArchivedAt:  src.ArchivedAt,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
	}
}

func mapFormMember(src sqlc.FormMember) models.FormMember {
	return models.FormMember{
		UserID:  src.UserID,
		FormID:  src.FormID,
		Role:    models.FormMemberRole(src.Role),
		AddedAt: src.AddedAt,
		AddedBy: src.AddedBy,
	}
}
