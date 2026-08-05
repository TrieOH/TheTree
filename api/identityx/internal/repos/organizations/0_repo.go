package organizations

import (
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/database"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.OrganizationRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("organization"),
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
