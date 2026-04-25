package projects

import (
	"Informd/internal/platform/database"
	"Informd/internal/platform/database/sqlc"
	"Informd/internal/shared/contracts"
	"Informd/internal/shared/errx"
	"Informd/internal/shared/ports"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type repo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ ports.ProjectsRepo = (*repo)(nil)

func NewProjectRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.ProjectsRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func (repo *repo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func mapProjectFromDB(src *sqlc.Project) *contracts.Project {
	return &contracts.Project{
		ID:        src.ID,
		OwnerID:   src.OwnerID,
		Name:      src.Name,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
	}
}

func (repo *repo) Create(ctx context.Context, toCreate contracts.Project) (*contracts.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectsRepo.Create")
	defer span.End()

	sqlcProject, err := repo.queries(ctx).CreateProject(ctx, sqlc.CreateProjectParams{
		ID:      toCreate.ID,
		OwnerID: toCreate.OwnerID,
		Name:    toCreate.Name,
	})
	if err != nil {
		return nil, errx.FromDB(err, "project")
	}

	return mapProjectFromDB(&sqlcProject), nil
}

func (repo *repo) GetByID(ctx context.Context, id uuid.UUID) (*contracts.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectsRepo.GetByID")
	defer span.End()

	sqlcProject, err := repo.queries(ctx).GetProjectByID(ctx, id)
	if err != nil {
		return nil, errx.FromDB(err, "project")
	}

	return mapProjectFromDB(&sqlcProject), nil
}
func (repo *repo) GetByName(ctx context.Context, name string, ownerID uuid.UUID) (*contracts.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectsRepo.GetByName")
	defer span.End()

	sqlcProject, err := repo.queries(ctx).GetProjectByName(ctx, sqlc.GetProjectByNameParams{
		OwnerID: ownerID,
		Name:    name,
	})
	if err != nil {
		return nil, errx.FromDB(err, "project")
	}

	return mapProjectFromDB(&sqlcProject), nil
}

func (repo *repo) List(ctx context.Context, ownerID uuid.UUID) ([]contracts.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectsRepo.List")
	defer span.End()

	sqlcProjects, err := repo.queries(ctx).ListProjectsByOwner(ctx, ownerID)
	if err != nil {
		return nil, errx.FromDB(err, "project")
	}

	out := make([]contracts.Project, 0, len(sqlcProjects))
	for _, project := range sqlcProjects {
		out = append(out, *mapProjectFromDB(&project))
	}
	return out, nil
}

func (repo *repo) ListByIDs(ctx context.Context, ids []string) ([]contracts.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectsRepo.ListByIDs")
	defer span.End()

	uuids := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		parsed, err := uuid.Parse(id)
		if err != nil {
			continue
		}
		uuids = append(uuids, parsed)
	}

	sqlcProjects, err := repo.queries(ctx).ListProjectsByIDs(ctx, uuids)
	if err != nil {
		return nil, errx.FromDB(err, "project")
	}

	out := make([]contracts.Project, 0, len(sqlcProjects))
	for _, project := range sqlcProjects {
		out = append(out, *mapProjectFromDB(&project))
	}
	return out, nil
}
