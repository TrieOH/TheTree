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

var _ ports.ProjectRepo = (*Repo)(nil)

func NewRepo(q *sqlc2.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("project"),
	}
}

func mapProject(src sqlc2.Project) models.Project {
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

func mapProjectMember(src sqlc2.ProjectMember) models.ProjectMember {
	return models.ProjectMember{
		ProjectID: src.ProjectID,
		ActorID:   src.ActorID,
		Role:      models.ProjectRole(src.Role),
		Metadata:  src.Metadata,
		JoinedAt:  src.JoinedAt,
	}
}
