package platform_roles

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

var _ ports.PlatformRolesRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("platform role"),
	}
}

func mapPlatformRole(src sqlc.PlatformRole) models.PlatformRoleRelation {
	return models.PlatformRoleRelation{
		ActorID:   src.ActorID,
		Role:      models.PlatformRole(src.Role),
		Metadata:  src.Metadata,
		CreatedAt: src.CreatedAt,
	}
}
