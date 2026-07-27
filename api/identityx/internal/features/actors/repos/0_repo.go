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

var _ ports.ActorRepo = (*Repo)(nil)

func NewRepo(q *sqlc2.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler("actor"),
	}
}

func mapActor(src sqlc2.Actor) models.Actor {
	return models.Actor{
		ID:           src.ID,
		ProjectID:    src.ProjectID,
		AuthMethod:   models.AuthMethod(src.AuthMethod),
		VerifiedAt:   src.VerifiedAt,
		PasswordHash: src.PasswordHash,
		Email:        src.Email,
		Type:         models.ActorType(src.Type),
		Metadata:     src.Metadata,
		CreatedAt:    src.CreatedAt,
		UpdatedAt:    src.UpdatedAt,
		DeletedAt:    src.DeletedAt,
	}
}
