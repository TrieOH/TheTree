package repo

import (
	"context"
	"lib/database"
	"payssage/errs"
	"payssage/internal/sqlc"
	"payssage/models"
	"payssage/ports"

	"github.com/google/uuid"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.OrganizationRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrOrganizationNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapOrganization(src sqlc.Organization) models.Organization {
	return models.Organization{
		ID:        src.ID,
		OwnerID:   src.OwnerID,
		Name:      src.Name,
		Slug:      src.Slug,
		CreatedAt: src.CreatedAt,
		DeletedAt: src.DeletedAt,
	}
}

func mapOrganizationMember(src sqlc.OrgMember) models.OrganizationMember {
	return models.OrganizationMember{
		OrganizationID: src.OrganizationID,
		MemberID:       src.MemberID,
		Role:           models.OrganizationRole(src.Role),
		JoinedAt:       src.JoinedAt,
	}
}

func (repo *Repo) AddMember(ctx context.Context, toAdd models.OrganizationMember) error {
	err := repo.db(ctx).AddOrganizationMember(ctx, sqlc.AddOrganizationMemberParams{
		MemberID:       toAdd.MemberID,
		OrganizationID: toAdd.OrganizationID,
		Role:           string(toAdd.Role),
	})
	return repo.dbe(err)
}

func (repo *Repo) Create(ctx context.Context, toCreate models.Organization) (*models.Organization, error) {
	sqlcOrg, err := repo.db(ctx).CreateOrganization(ctx, sqlc.CreateOrganizationParams{
		OwnerID: toCreate.OwnerID,
		Name:    toCreate.Name,
		Slug:    toCreate.Slug,
	})
	return database.One(sqlcOrg, err, repo.dbe, mapOrganization)
}

func (repo *Repo) GetByID(ctx context.Context, orgID uuid.UUID) (*models.Organization, error) {
	sqlcOrg, err := repo.db(ctx).GetOrganizationByID(ctx, orgID)
	return database.One(sqlcOrg, err, repo.dbe, mapOrganization)
}

func (repo *Repo) GetMember(ctx context.Context, memberID, orgID uuid.UUID) (*models.OrganizationMember, error) {
	sqlcMember, err := repo.db(ctx).GetOrganizationMember(ctx, sqlc.GetOrganizationMemberParams{
		MemberID:       memberID,
		OrganizationID: orgID,
	})
	return database.One(sqlcMember, err, repo.dbe, mapOrganizationMember)
}

func (repo *Repo) GetRole(ctx context.Context, actorID, orgID uuid.UUID) (models.OrganizationRole, error) {
	member, err := repo.GetMember(ctx, actorID, orgID)
	if err != nil {
		return "", err
	}
	return member.Role, nil
}

func (repo *Repo) ListJoined(ctx context.Context, userID uuid.UUID) ([]models.Organization, error) {
	sqlcOrgs, err := repo.db(ctx).ListJoinedOrganizations(ctx, userID)
	return database.Many(sqlcOrgs, err, repo.dbe, mapOrganization)
}

func (repo *Repo) ListMembers(ctx context.Context, orgID uuid.UUID) ([]models.OrganizationMember, error) {
	sqlcMembers, err := repo.db(ctx).ListOrganizationMembers(ctx, orgID)
	return database.Many(sqlcMembers, err, repo.dbe, mapOrganizationMember)
}

func (repo *Repo) ListOwned(ctx context.Context, userID uuid.UUID) ([]models.Organization, error) {
	sqlcOrgs, err := repo.db(ctx).ListOwnedOrganizations(ctx, userID)
	return database.Many(sqlcOrgs, err, repo.dbe, mapOrganization)
}

func (repo *Repo) RemoveMember(ctx context.Context, memberID, orgID uuid.UUID) error {
	err := repo.db(ctx).RemoveOrganizationMember(ctx, sqlc.RemoveOrganizationMemberParams{
		MemberID:       memberID,
		OrganizationID: orgID,
	})
	return repo.dbe(err)
}
