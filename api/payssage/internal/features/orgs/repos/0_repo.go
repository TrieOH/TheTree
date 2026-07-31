package repos

import (
	"lib/database"
	"payssage/internal/sqlc"
	"payssage/models"
	"payssage/ports"

	"go.opentelemetry.io/otel/trace"
)

type Repo struct {
	q      *sqlc.Queries
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.OrganizationRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries, tracer trace.Tracer) *Repo {
	return &Repo{
		q:      q,
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
