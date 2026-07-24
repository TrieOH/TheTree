package repos

import (
	"lib/database"
	"univents/internal/database/sqlc"
	"univents/models"
	"univents/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type repo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.CertificationRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.CertificationRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("certifications"),
	}
}

func mapCertificationTemplate(src sqlc.CertificationTemplate) models.CertificationTemplate {
	return models.CertificationTemplate{
		ID:        src.ID,
		EditionID: src.EditionID,
		Title:     src.Title,
		Data:      src.Data,
		URL:       src.Url,
		CreatedAt: src.CreatedAt,
	}
}

func mapCertification(src sqlc.Certification) models.Certification {
	return models.Certification{
		ID:          src.ID,
		UserID:      src.UserID,
		TargetID:    src.TargetID,
		TargetType:  src.TargetType,
		CertifiedAt: src.CertifiedAt,
		Hash:        src.Hash,
	}
}
