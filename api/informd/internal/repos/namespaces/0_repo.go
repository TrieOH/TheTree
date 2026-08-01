package namespaces

import (
	"Informd/internal/sqlc"
	"Informd/models"
	"Informd/ports"
	"lib/database"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.NamespaceRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("namespace"),
	}
}

func mapNamespace(src sqlc.Namespace) models.Namespace {
	return models.Namespace{
		ID:        src.ID,
		OwnerID:   src.OwnerID,
		Name:      src.Name,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
	}
}

func mapNamespaceMember(src sqlc.NamespaceMember) models.NamespaceMember {
	return models.NamespaceMember{
		UserID:      src.UserID,
		NamespaceID: src.NamespaceID,
		Role:        models.NamespaceMemberRole(src.Role),
		AddedAt:     src.AddedAt,
		AddedBy:     src.AddedBy,
	}
}
