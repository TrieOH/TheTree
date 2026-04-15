package projects

import (
	"IdentityX/internal/platform/database"
	sqlc2 "IdentityX/internal/platform/database/sqlc"
	"IdentityX/internal/shared/contracts"
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/ports"
	"context"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type projectRepo struct {
	q      *sqlc2.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
}

var _ ports.ProjectRepository = (*projectRepo)(nil)

func NewRepo(q *sqlc2.Queries, log *zap.Logger, tracer trace.Tracer) ports.ProjectRepository {
	return &projectRepo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func (repo *projectRepo) queries(ctx context.Context) *sqlc2.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func mapProjectFromDB(dst *contracts.Project, src *sqlc2.Project) {
	dst.ID = src.ID
	dst.ProjectName = src.ProjectName
	dst.Domain = src.Domain
	dst.OwnerID = src.OwnerID
	dst.Metadata = src.Metadata
	dst.IsActive = src.IsActive
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
}

func (repo *projectRepo) Create(ctx context.Context, toCreate contracts.Project) (*contracts.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectRepo.Create",
		trace.WithAttributes(
			attribute.String("project.owner_id", toCreate.OwnerID.String()),
			attribute.String("project.name", toCreate.ProjectName),
		),
	)
	defer span.End()

	sqlcProject, err := repo.queries(ctx).CreateProject(ctx, sqlc2.CreateProjectParams{
		ProjectName: toCreate.ProjectName,
		Domain:      toCreate.Domain,
		OwnerID:     toCreate.OwnerID,
		Metadata:    toCreate.Metadata,
		IsActive:    toCreate.IsActive,
	})
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(attribute.String("project.id", sqlcProject.ID.String()))

	mapProjectFromDB(&toCreate, &sqlcProject)
	return &toCreate, nil
}

func (repo *projectRepo) GetByIDExternal(ctx context.Context, projectID, ownerID uuid.UUID) (*contracts.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectRepo.GetByIDExternal",
		trace.WithAttributes(
			attribute.String("project.owner_id", ownerID.String()),
			attribute.String("project.id", projectID.String()),
		),
	)
	defer span.End()

	sqlcProject, err := repo.queries(ctx).GetProjectByIDExternal(ctx, sqlc2.GetProjectByIDExternalParams{
		ID:      projectID,
		OwnerID: ownerID,
	})
	if err != nil {
		return nil, fail.From(err).WithArgs("project").RecordCtx(ctx)
	}

	span.SetAttributes(attribute.String("project.name", sqlcProject.ProjectName))

	var proj contracts.Project
	mapProjectFromDB(&proj, &sqlcProject)
	return &proj, nil
}

func (repo *projectRepo) GetByIDInternal(ctx context.Context, projectID uuid.UUID) (*contracts.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectRepo.GetByIDInternal",
		trace.WithAttributes(
			attribute.String("project.id", projectID.String()),
		),
	)
	defer span.End()

	sqlcProject, err := repo.queries(ctx).GetProjectByIDInternal(ctx, projectID)
	if err != nil {
		return nil, fail.From(err).WithArgs("project").RecordCtx(ctx)
	}

	span.SetAttributes(attribute.String("project.name", sqlcProject.ProjectName))

	var proj contracts.Project
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

	isOwner, err := repo.queries(ctx).IsOwnerOf(ctx, sqlc2.IsOwnerOfParams{
		OwnerID: ownerID,
		ID:      projectID,
	})
	if err != nil {
		return false, fail.From(err).RecordCtx(ctx)
	}

	return isOwner, nil
}

func (repo *projectRepo) List(ctx context.Context, ownerID uuid.UUID) ([]contracts.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectRepo.List",
		trace.WithAttributes(
			attribute.String("project.owner_id", ownerID.String()),
		),
	)
	defer span.End()

	sqlcProjects, err := repo.queries(ctx).ListProjects(ctx, ownerID)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(attribute.Int("project.count", len(sqlcProjects)))

	projects := make([]contracts.Project, 0, len(sqlcProjects))
	for _, sqlcProject := range sqlcProjects {
		var proj contracts.Project
		mapProjectFromDB(&proj, &sqlcProject)
		projects = append(projects, proj)
	}

	return projects, nil
}

func (repo *projectRepo) Update(ctx context.Context, toUpdate contracts.Project, ownerID uuid.UUID) (*contracts.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "ProjectRepo.Update",
		trace.WithAttributes(
			attribute.String("project.owner_id", ownerID.String()),
			attribute.String("project.id", toUpdate.ID.String()),
		),
	)
	defer span.End()

	sqlcProject, err := repo.queries(ctx).UpdateProject(ctx, sqlc2.UpdateProjectParams{
		ID:          toUpdate.ID,
		OwnerID:     ownerID,
		ProjectName: toUpdate.ProjectName,
		Domain:      toUpdate.Domain,
		Metadata:    toUpdate.Metadata,
	})
	if err != nil {
		return nil, fail.From(err).WithArgs("project").RecordCtx(ctx)
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

	affectedRows, err := repo.queries(ctx).DeleteProject(ctx, sqlc2.DeleteProjectParams{
		ID:      projectID,
		OwnerID: ownerID,
	})

	if err != nil {
		return fail.From(err).WithArgs("project").RecordCtx(ctx)
	}

	if affectedRows == 0 {
		return fail.New(errx.ProjectNotFound)
	}

	return nil
}
