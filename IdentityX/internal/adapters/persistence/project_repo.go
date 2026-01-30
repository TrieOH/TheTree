package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/project"
	"GoAuth/internal/ports/outbounds"
	"context"

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
}

var _ outbounds.ProjectRepository = (*projectRepo)(nil)

func NewProjectRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) outbounds.ProjectRepository {
	return &projectRepo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func (repo *projectRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(transactions.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func mapProjectFromDB(dst *project.Project, src *sqlc.Project) {
	dst.ID = src.ID
	dst.ProjectName = src.ProjectName
	dst.OwnerID = src.OwnerID
	dst.Metadata = src.Metadata
	dst.IsActive = src.IsActive
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
}

func (repo *projectRepo) Create(ctx context.Context, toCreate project.Project) (*project.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectRepo.Create",
		trace.WithAttributes(
			attribute.String("project.owner_id", toCreate.OwnerID.String()),
			attribute.String("project.name", toCreate.ProjectName),
		),
	)
	defer span.End()

	sqlcProject, err := repo.queries(ctx).CreateProject(ctx, sqlc.CreateProjectParams{
		ProjectName: toCreate.ProjectName,
		OwnerID:     toCreate.OwnerID,
		Metadata:    toCreate.Metadata,
		IsActive:    toCreate.IsActive,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(attribute.String("project.id", sqlcProject.ID.String()))

	mapProjectFromDB(&toCreate, &sqlcProject)
	return &toCreate, nil
}

func (repo *projectRepo) GetByIDExternal(ctx context.Context, projectID, ownerID uuid.UUID) (*project.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectRepo.GetByIDExternal",
		trace.WithAttributes(
			attribute.String("project.owner_id", ownerID.String()),
			attribute.String("project.id", projectID.String()),
		),
	)
	defer span.End()

	sqlcProject, err := repo.queries(ctx).GetProjectByIDExternal(ctx, sqlc.GetProjectByIDExternalParams{
		ID:      projectID,
		OwnerID: ownerID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(attribute.String("project.name", sqlcProject.ProjectName))

	var proj project.Project
	mapProjectFromDB(&proj, &sqlcProject)
	return &proj, nil
}

func (repo *projectRepo) GetByIDInternal(ctx context.Context, projectID uuid.UUID) (*project.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectRepo.GetByIDInternal",
		trace.WithAttributes(
			attribute.String("project.id", projectID.String()),
		),
	)
	defer span.End()

	sqlcProject, err := repo.queries(ctx).GetProjectByIDInternal(ctx, projectID)
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(attribute.String("project.name", sqlcProject.ProjectName))

	var proj project.Project
	mapProjectFromDB(&proj, &sqlcProject)
	return &proj, nil
}

func (repo *projectRepo) IsOwnerOf(ctx context.Context, projectID, ownerID uuid.UUID) (bool, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectRepo.IsOwnerOf",
		trace.WithAttributes(
			attribute.String("project.owner_id", ownerID.String()),
			attribute.String("project.id", projectID.String()),
		),
	)
	defer span.End()

	isOwner, err := repo.queries(ctx).IsOwnerOf(ctx, sqlc.IsOwnerOfParams{
		OwnerID: ownerID,
		ID:      projectID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return false, sqlcErr
	}

	return isOwner, nil
}

func (repo *projectRepo) List(ctx context.Context, ownerID uuid.UUID) ([]project.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectRepo.List",
		trace.WithAttributes(
			attribute.String("project.owner_id", ownerID.String()),
		),
	)
	defer span.End()

	sqlcProjects, err := repo.queries(ctx).ListProjects(ctx, ownerID)
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(attribute.Int("project.count", len(sqlcProjects)))

	projects := make([]project.Project, 0, len(sqlcProjects))
	for _, sqlcProject := range sqlcProjects {
		var proj project.Project
		mapProjectFromDB(&proj, &sqlcProject)
		projects = append(projects, proj)
	}

	return projects, nil
}

func (repo *projectRepo) Update(ctx context.Context, toUpdate project.Project, ownerID uuid.UUID) (*project.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectRepo.Update",
		trace.WithAttributes(
			attribute.String("project.owner_id", ownerID.String()),
			attribute.String("project.id", toUpdate.ID.String()),
		),
	)
	defer span.End()

	sqlcProject, err := repo.queries(ctx).UpdateProject(ctx, sqlc.UpdateProjectParams{
		ID:          toUpdate.ID,
		OwnerID:     ownerID,
		ProjectName: toUpdate.ProjectName,
		Metadata:    toUpdate.Metadata,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	mapProjectFromDB(&toUpdate, &sqlcProject)
	return &toUpdate, nil
}

func (repo *projectRepo) Delete(ctx context.Context, projectID, ownerID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "ProjectRepo.Delete",
		trace.WithAttributes(
			attribute.String("project.owner_id", ownerID.String()),
			attribute.String("project.id", projectID.String()),
		),
	)
	defer span.End()

	affectedRows, err := repo.queries(ctx).DeleteProject(ctx, sqlc.DeleteProjectParams{
		ID:      projectID,
		OwnerID: ownerID,
	})

	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}

	// FIXME make me a generic error
	if affectedRows == 0 {
		return apierr.ErrNotFound.WithMsg("project not found").WithID(apierr.ProjectNotFound)
	}

	return nil
}
