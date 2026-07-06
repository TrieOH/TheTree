package repos

import (
	"lib/database"
	"payssage/internal/database/sqlc"
	"payssage/models"
	"payssage/ports"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type repo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.OrganizationRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.OrganizationRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("organization"),
	}
}

func mapOrganization(src sqlc.Organization) models.Organization {
	return models.Organization{
		ID:        src.ID,
		OwnerID:   src.OwnerID,
		Name:      src.Name,
		Slug:      src.Slug,
		CreatedAt: src.CreatedAt,
		DeletedAt: src.DeletedAt,
	}
}

func mapOrganizationMember(src sqlc.OrgMember) models.OrganizationMember {
	return models.OrganizationMember{
		OrganizationID: src.OrganizationID,
		MemberID:       src.MemberID,
		Role:           models.OrganizationRole(src.Role),
		JoinedAt:       src.JoinedAt,
	}
}
