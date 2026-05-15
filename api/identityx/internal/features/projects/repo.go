package projects

import (
	"IdentityX/contracts"
	"IdentityX/internal/platform/database/sqlc"
	"IdentityX/internal/shared/ports"
	"context"
	"lib/database"
	"lib/errx"
	"lib/xslices"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type projectRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
	dbe    *errx.DBHandler
}

var _ ports.ProjectRepository = (*projectRepo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer, dbe *errx.DBHandler) ports.ProjectRepository {
	return &projectRepo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    dbe,
	}
}

func (repo *projectRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func (repo *projectRepo) span(ctx context.Context, op string) (context.Context, trace.Span) {
	return repo.tracer.Start(ctx, "ProjectRepo."+op)
}

func mapProjectFromDB(src sqlc.Project) contracts.Project {
	return contracts.Project{
		ID:          src.ID,
		ProjectName: src.ProjectName,
		Domain:      src.Domain,
		OwnerID:     src.OwnerID,
		Metadata:    src.Metadata,
		IsActive:    src.IsActive,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
	}
}

func (repo *projectRepo) Create(ctx context.Context, toCreate contracts.Project) (*contracts.Project, error) {
	ctx, span := repo.span(ctx, "Create")
	span.SetAttributes(attribute.String("project.owner_id", toCreate.OwnerID.String()))
	span.SetAttributes(attribute.String("project.name", toCreate.ProjectName))
	defer span.End()
	sqlcProject, err := repo.queries(ctx).CreateProject(ctx, sqlc.CreateProjectParams{
		ProjectName: toCreate.ProjectName,
		Domain:      toCreate.Domain,
		OwnerID:     toCreate.OwnerID,
		Metadata:    toCreate.Metadata,
		IsActive:    toCreate.IsActive,
	})
	if err != nil {
		return nil, repo.dbe.DB(err, "project")
	}
	span.SetAttributes(attribute.String("project.id", sqlcProject.ID.String()))
	return new(mapProjectFromDB(sqlcProject)), nil
}

func (repo *projectRepo) GetByIDExternal(ctx context.Context, projectID, ownerID uuid.UUID) (*contracts.Project, error) {
	ctx, span := repo.span(ctx, "GetByIDExternal")
	span.SetAttributes(attribute.String("project.owner_id", ownerID.String()))
	span.SetAttributes(attribute.String("project.id", projectID.String()))
	defer span.End()
	sqlcProject, err := repo.queries(ctx).GetProjectByIDExternal(ctx, sqlc.GetProjectByIDExternalParams{
		ID:      projectID,
		OwnerID: ownerID,
	})
	if err != nil {
		return nil, repo.dbe.DB(err, "project")
	}
	span.SetAttributes(attribute.String("project.name", sqlcProject.ProjectName))
	return new(mapProjectFromDB(sqlcProject)), nil
}

func (repo *projectRepo) GetByIDInternal(ctx context.Context, projectID uuid.UUID) (*contracts.Project, error) {
	ctx, span := repo.span(ctx, "GetByIDInternal")
	span.SetAttributes(attribute.String("project.id", projectID.String()))
	defer span.End()
	sqlcProject, err := repo.queries(ctx).GetProjectByIDInternal(ctx, projectID)
	if err != nil {
		return nil, repo.dbe.DB(err, "project")
	}
	span.SetAttributes(attribute.String("project.name", sqlcProject.ProjectName))
	return new(mapProjectFromDB(sqlcProject)), nil
}

func (repo *projectRepo) IsOwnerOf(ctx context.Context, projectID, ownerID uuid.UUID) (bool, error) {
	ctx, span := repo.span(ctx, "IsOwnerOf")
	span.SetAttributes(attribute.String("project.owner_id", ownerID.String()))
	span.SetAttributes(attribute.String("project.id", projectID.String()))
	defer span.End()
	isOwner, err := repo.queries(ctx).IsOwnerOf(ctx, sqlc.IsOwnerOfParams{
		OwnerID: ownerID,
		ID:      projectID,
	})
	if err != nil {
		return false, repo.dbe.DB(err, "project")
	}
	return isOwner, nil
}

func (repo *projectRepo) List(ctx context.Context, ownerID uuid.UUID) ([]contracts.Project, error) {
	ctx, span := repo.span(ctx, "List")
	span.SetAttributes(attribute.String("project.owner_id", ownerID.String()))
	defer span.End()
	sqlcProjects, err := repo.queries(ctx).ListProjects(ctx, ownerID)
	if err != nil {
		return nil, repo.dbe.DB(err, "project")
	}
	span.SetAttributes(attribute.Int("project.count", len(sqlcProjects)))
	return xslices.MapSlice(sqlcProjects, mapProjectFromDB), nil
}

func (repo *projectRepo) Update(ctx context.Context, toUpdate contracts.Project, ownerID uuid.UUID) (*contracts.Project, error) {
	ctx, span := repo.span(ctx, "Update")
	span.SetAttributes(attribute.String("project.owner_id", ownerID.String()))
	span.SetAttributes(attribute.String("project.id", toUpdate.ID.String()))
	defer span.End()
	sqlcProject, err := repo.queries(ctx).UpdateProject(ctx, sqlc.UpdateProjectParams{
		ID:          toUpdate.ID,
		OwnerID:     ownerID,
		ProjectName: toUpdate.ProjectName,
		Domain:      toUpdate.Domain,
		Metadata:    toUpdate.Metadata,
	})
	if err != nil {
		return nil, repo.dbe.DB(err, "project")
	}
	return new(mapProjectFromDB(sqlcProject)), nil
}

func (repo *projectRepo) Delete(ctx context.Context, projectID, ownerID uuid.UUID) error {
	ctx, span := repo.span(ctx, "Delete")
	span.SetAttributes(attribute.String("project.owner_id", ownerID.String()))
	span.SetAttributes(attribute.String("project.id", projectID.String()))
	defer span.End()
	affectedRows, err := repo.queries(ctx).DeleteProject(ctx, sqlc.DeleteProjectParams{
		ID:      projectID,
		OwnerID: ownerID,
	})
	if err != nil {
		return repo.dbe.DB(err, "project")
	}
	if affectedRows == 0 {
		return fun.ErrNotFound("project not found")
	}
	return nil
}
