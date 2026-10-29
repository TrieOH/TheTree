package projects

import (
	"IdentityX/errs"
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"context"
	"lib/database"

	"github.com/google/uuid"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.ProjectRepo = (*Repo)(nil)

func NewRepo(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrProjectNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapProject(src sqlc.Project) models.Project {
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

func mapProjectMember(src sqlc.ProjectMember) models.ProjectMember {
	return models.ProjectMember{
		ProjectID: src.ProjectID,
		ActorID:   src.ActorID,
		Role:      models.ProjectRole(src.Role),
		Metadata:  src.Metadata,
		JoinedAt:  src.JoinedAt,
	}
}

func (repo *Repo) AddMember(ctx context.Context, toCreate models.ProjectMember) error {
	err := repo.db(ctx).AddProjectMember(ctx, sqlc.AddProjectMemberParams{
		ProjectID: toCreate.ProjectID,
		ActorID:   toCreate.ActorID,
		Role:      string(toCreate.Role),
		Metadata:  toCreate.Metadata,
	})
	return repo.dbe(err)
}

func (repo *Repo) Create(ctx context.Context, project models.Project) (*models.Project, error) {
	row, err := repo.db(ctx).CreateProject(ctx, sqlc.CreateProjectParams{
		OrganizationID: project.OrganizationID,
		OwnerID:        project.OwnerID,
		Name:           project.Name,
		Domain:         project.Domain,
		BrandSlug:      project.BrandSlug,
		Metadata:       project.Metadata,
	})
	return database.One(row, err, repo.dbe, mapProject)
}

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	sqlcProject, err := repo.db(ctx).GetProjectByID(ctx, id)
	return database.One(sqlcProject, err, repo.dbe, mapProject)
}

func (repo *Repo) GetMember(ctx context.Context, actorID, projectID uuid.UUID) (*models.ProjectMember, error) {
	sqlcMember, err := repo.db(ctx).GetProjectMemberByID(ctx, sqlc.GetProjectMemberByIDParams{
		ActorID:   actorID,
		ProjectID: projectID,
	})
	return database.One(sqlcMember, err, repo.dbe, mapProjectMember)
}

func (repo *Repo) GetRole(ctx context.Context, actorID, projectID uuid.UUID) (models.ProjectRole, error) {
	member, err := repo.GetMember(ctx, actorID, projectID)
	if err != nil {
		return "", err
	}
	return member.Role, nil
}

// ListAll returns every project — the Key-lifecycle module crosses it so
// EnsureAll can provision every scope (platform + each project).
func (repo *Repo) ListAll(ctx context.Context) ([]models.Project, error) {
	rows, err := repo.db(ctx).ListProjects(ctx)
	if err != nil {
		return nil, repo.dbe(err)
	}
	out := make([]models.Project, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapProject(row))
	}
	return out, nil
}

func (repo *Repo) ListByOrganization(ctx context.Context, orgID uuid.UUID) ([]models.Project, error) {
	sqlcProjects, err := repo.db(ctx).ListProjectsByOrganizationID(ctx, &orgID)
	return database.Many(sqlcProjects, err, repo.dbe, mapProject)
}

func (repo *Repo) ListJoined(ctx context.Context, userID uuid.UUID) ([]models.Project, error) {
	sqlcProjects, err := repo.db(ctx).ListJoinedProjects(ctx, userID)
	return database.Many(sqlcProjects, err, repo.dbe, mapProject)
}

func (repo *Repo) ListMembers(ctx context.Context, projectID uuid.UUID) ([]models.ProjectMember, error) {
	sqlcMembers, err := repo.db(ctx).ListProjectMembers(ctx, projectID)
	return database.Many(sqlcMembers, err, repo.dbe, mapProjectMember)
}

func (repo *Repo) ListOwned(ctx context.Context, userID uuid.UUID) ([]models.Project, error) {
	sqlcProjects, err := repo.db(ctx).ListOwnedProjects(ctx, userID)
	return database.Many(sqlcProjects, err, repo.dbe, mapProject)
}

func (repo *Repo) RemoveMember(ctx context.Context, actorID, projectID uuid.UUID) error {
	err := repo.db(ctx).RemoveProjectMember(ctx, sqlc.RemoveProjectMemberParams{
		ActorID:   actorID,
		ProjectID: projectID,
	})
	return repo.dbe(err)
}
