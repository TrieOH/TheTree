package repo

import (
	"IdentityX/errs"
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"context"
	"encoding/json"
	"lib/database"

	"github.com/google/uuid"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.PlatformRolesRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrPlatformRoleNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapPlatformRole(src sqlc.PlatformRole) models.PlatformRoleRelation {
	return models.PlatformRoleRelation{
		ActorID:   src.ActorID,
		Role:      models.PlatformRole(src.Role),
		Metadata:  src.Metadata,
		CreatedAt: src.CreatedAt,
	}
}

// GetRole returns the actor's platform role. A NotFound surfaces for actors
// with no platform role — callers map it (authz.CheckPlatform forbids).
func (repo *Repo) GetRole(ctx context.Context, actorID uuid.UUID) (models.PlatformRole, error) {
	role, err := repo.db(ctx).GetPlatformRoleByActor(ctx, actorID)
	if err != nil {
		return "", repo.dbe(err)
	}
	return models.PlatformRole(role), nil
}

func (repo *Repo) Give(ctx context.Context, actorID uuid.UUID, role models.PlatformRole, metadata *json.RawMessage) (*models.PlatformRoleRelation, error) {
	sqlcPlatformRole, err := repo.db(ctx).GivePlatformRole(ctx, sqlc.GivePlatformRoleParams{
		ActorID:  actorID,
		Role:     string(role),
		Metadata: metadata,
	})
	return database.One(sqlcPlatformRole, err, repo.dbe, mapPlatformRole)
}
