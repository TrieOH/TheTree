package repos

import (
	sqlc2 "IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/database"
)

type Repo struct {
	q   *sqlc2.Queries
	dbe database.ErrorHandler
}

var _ ports.OrganizationRepo = (*Repo)(nil)

func NewRepo(q *sqlc2.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("organization"),
	}
}

func mapOrganization(src sqlc2.Organization) models.Organization {
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

func mapOrganizationMember(src sqlc2.OrgMember) models.OrganizationMember {
	return models.OrganizationMember{
		OrganizationID: src.OrganizationID,
		ActorID:        src.ActorID,
		Role:           models.OrganizationRole(src.Role),
		Metadata:       src.Metadata,
		JoinedAt:       src.JoinedAt,
	}
}
