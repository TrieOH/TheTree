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

var _ ports.ProjectRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.ProjectRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("project"),
	}
}

func mapProject(src sqlc.Project) models.Project {
	return models.Project{
		ID:               src.ID,
		OrganizationID:   src.OrganizationID,
		OwnerID:          src.OwnerID,
		Name:             src.Name,
		Domain:           src.Domain,
		BrandSlug:        src.BrandSlug,
		DomainVerifiedAt: src.DomainVerifiedAt,
		Metadata:         src.Metadata,
		CreatedAt:        src.CreatedAt,
		DeletedAt:        src.DeletedAt,
	}
}

func mapProjectMember(src sqlc.ProjectMember) models.ProjectMember {
	return models.ProjectMember{
		ProjectID: src.ProjectID,
		ActorID:   src.ActorID,
		Role:      models.ProjectRole(src.Role),
		Metadata:  src.Metadata,
		JoinedAt:  src.JoinedAt,
	}
}
