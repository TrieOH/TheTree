package repo

import (
	"IdentityX/errs"
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"context"
	"lib/database"
	"time"

	"github.com/google/uuid"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.ActorRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrActorNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapActor(src sqlc.Actor) models.Actor {
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

func (repo *Repo) GetByEmail(ctx context.Context, email string, projectID *uuid.UUID) (*models.Actor, error) {
	sqlcActor, err := repo.db(ctx).GetActorByEmail(ctx, sqlc.GetActorByEmailParams{
		Email:     &email,
		ProjectID: projectID,
	})
	return database.One(sqlcActor, err, repo.dbe, mapActor)
}

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Actor, error) {
	sqlcActor, err := repo.db(ctx).GetActorByID(ctx, id)
	return database.One(sqlcActor, err, repo.dbe, mapActor)
}

func (repo *Repo) GetProjectServiceAccount(ctx context.Context, id uuid.UUID) (*models.Actor, error) {
	sqlcActor, err := repo.db(ctx).GetProjectServiceAccount(ctx, &id)
	return database.One(sqlcActor, err, repo.dbe, mapActor)
}

func (repo *Repo) List(ctx context.Context, projectID uuid.UUID) ([]models.Actor, error) {
	sqlcActors, err := repo.db(ctx).ListActorsFromProject(ctx, &projectID)
	return database.Many(sqlcActors, err, repo.dbe, mapActor)
}

func (repo *Repo) Register(ctx context.Context, toRegister models.Actor) (*models.Actor, error) {
	sqlcActor, err := repo.db(ctx).RegisterActor(ctx, sqlc.RegisterActorParams{
		ProjectID:    toRegister.ProjectID,
		AuthMethod:   string(toRegister.AuthMethod),
		PasswordHash: toRegister.PasswordHash,
		Email:        toRegister.Email,
		Type:         string(toRegister.Type),
		Metadata:     toRegister.Metadata,
		VerifiedAt:   toRegister.VerifiedAt,
	})
	return database.One(sqlcActor, err, repo.dbe, mapActor)
}

func (repo *Repo) UpdateLastLoginAt(ctx context.Context, actorID uuid.UUID) error {
	return repo.dbe(repo.db(ctx).UpdateActorLastLoginAt(ctx, actorID))
}

func (repo *Repo) SetVerifiedAt(ctx context.Context, actorID uuid.UUID, at time.Time) error {
	verifiedAt := at
	err := repo.db(ctx).UpdateActorVerifiedAt(ctx, sqlc.UpdateActorVerifiedAtParams{
		ActorID:    actorID,
		VerifiedAt: &verifiedAt,
	})
	return repo.dbe(err)
}

func (repo *Repo) UpdatePasswordHash(ctx context.Context, actorID uuid.UUID, hash string) error {
	passwordHash := hash
	err := repo.db(ctx).UpdateActorPasswordHash(ctx, sqlc.UpdateActorPasswordHashParams{
		ActorID:      actorID,
		PasswordHash: &passwordHash,
	})
	return repo.dbe(err)
}
