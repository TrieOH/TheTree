package repos

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/database"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Repo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.OrganizationRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) *Repo {
	return &Repo{
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
		Metadata:  src.Metadata,
		CreatedAt: src.CreatedAt,
		DeletedAt: src.DeletedAt,
	}
}

func mapOrganizationMember(src sqlc.OrgMember) models.OrganizationMember {
	return models.OrganizationMember{
		OrganizationID: src.OrganizationID,
		ActorID:        src.ActorID,
		Role:           models.OrganizationRole(src.Role),
		Metadata:       src.Metadata,
		JoinedAt:       src.JoinedAt,
	}
}
