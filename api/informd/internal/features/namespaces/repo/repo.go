package repo

import (
	"Informd/errs"
	"Informd/internal/sqlc"
	"Informd/models"
	"Informd/ports"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.NamespaceRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrNamespaceNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

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

func (repo *Repo) AddMember(ctx context.Context, toCreate models.NamespaceMember) error {
	err := repo.db(ctx).AddNamespaceMember(ctx, sqlc.AddNamespaceMemberParams{
		UserID:      toCreate.UserID,
		NamespaceID: toCreate.NamespaceID,
		Role:        string(toCreate.Role),
		AddedAt:     toCreate.AddedAt,
		AddedBy:     toCreate.AddedBy,
	})
	return repo.dbe(err)
}

func (repo *Repo) Create(ctx context.Context, toCreate models.Namespace) (*models.Namespace, error) {
	sqlcProject, err := repo.db(ctx).CreateNamespace(ctx, sqlc.CreateNamespaceParams{
		OwnerID: toCreate.OwnerID,
		Name:    toCreate.Name,
	})
	return database.One(sqlcProject, err, repo.dbe, mapNamespace)
}

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Namespace, error) {
	sqlcProject, err := repo.db(ctx).GetNamespaceByID(ctx, id)
	return database.One(sqlcProject, err, repo.dbe, mapNamespace)
}

func (repo *Repo) GetByName(ctx context.Context, name string, ownerID uuid.UUID) (*models.Namespace, error) {
	sqlcProject, err := repo.db(ctx).GetNamespaceByName(ctx, sqlc.GetNamespaceByNameParams{
		OwnerID: ownerID,
		Name:    name,
	})
	return database.One(sqlcProject, err, repo.dbe, mapNamespace)
}

func (repo *Repo) GetMember(ctx context.Context, userID, namespaceID uuid.UUID) (*models.NamespaceMember, error) {
	sqlcMember, err := repo.db(ctx).GetNamespaceMember(ctx, sqlc.GetNamespaceMemberParams{
		UserID:      userID,
		NamespaceID: namespaceID,
	})
	return database.One(sqlcMember, err, repo.dbe, mapNamespaceMember)
}

func (repo *Repo) ListJoined(ctx context.Context, userID uuid.UUID) ([]models.Namespace, error) {
	sqlcNamespaces, err := repo.db(ctx).ListJoinedNamespaces(ctx, userID)
	return database.Many(sqlcNamespaces, err, repo.dbe, mapNamespace)
}

func (repo *Repo) ListMembers(ctx context.Context, namespaceID uuid.UUID) ([]models.NamespaceMember, error) {
	sqlcMembers, err := repo.db(ctx).ListNamespaceMembers(ctx, namespaceID)
	return database.Many(sqlcMembers, err, repo.dbe, mapNamespaceMember)
}

func (repo *Repo) ListOwned(ctx context.Context, userID uuid.UUID) ([]models.Namespace, error) {
	sqlcNamespaces, err := repo.db(ctx).ListOwnedNamespaces(ctx, userID)
	return database.Many(sqlcNamespaces, err, repo.dbe, mapNamespace)
}

func (repo *Repo) RemoveMember(ctx context.Context, userID, namespaceID uuid.UUID) error {
	err := repo.db(ctx).RemoveNamespaceMember(ctx, sqlc.RemoveNamespaceMemberParams{
		UserID:      userID,
		NamespaceID: namespaceID,
	})
	return repo.dbe(err)
}

func (repo *Repo) GetRole(ctx context.Context, actorID, namespaceID uuid.UUID) (models.NamespaceMemberRole, error) {
	member, err := repo.GetMember(ctx, actorID, namespaceID)
	if err != nil {
		return "", err
	}
	return member.Role, nil
}
