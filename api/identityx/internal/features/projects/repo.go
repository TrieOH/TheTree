package projects

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/internal/shared/ports"
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type projectRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.ProjectRepository = (*projectRepo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.ProjectRepository {
	return &projectRepo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("project"),
	}
}

func mapProjectFromDB(src sqlc.Project) models.Project {
	return models.Project{
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

func (repo *projectRepo) Create(ctx context.Context, toCreate models.Project) (*models.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "Create")
	span.SetAttributes(attribute.String("project.owner_id", toCreate.OwnerID.String()))
	span.SetAttributes(attribute.String("project.name", toCreate.ProjectName))
	defer span.End()
	sqlcProject, err := database.Queries(ctx, repo.q).CreateProject(ctx, sqlc.CreateProjectParams{
		ProjectName: toCreate.ProjectName,
		Domain:      toCreate.Domain,
		OwnerID:     toCreate.OwnerID,
		Metadata:    toCreate.Metadata,
		IsActive:    toCreate.IsActive,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	span.SetAttributes(attribute.String("project.id", sqlcProject.ID.String()))
	return new(mapProjectFromDB(sqlcProject)), nil
}

func (repo *projectRepo) GetByIDExternal(ctx context.Context, projectID, ownerID uuid.UUID) (*models.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "GetByIDExternal")
	span.SetAttributes(attribute.String("project.owner_id", ownerID.String()))
	span.SetAttributes(attribute.String("project.id", projectID.String()))
	defer span.End()
	sqlcProject, err := database.Queries(ctx, repo.q).GetProjectByIDExternal(ctx, sqlc.GetProjectByIDExternalParams{
		ID:      projectID,
		OwnerID: ownerID,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	span.SetAttributes(attribute.String("project.name", sqlcProject.ProjectName))
	return new(mapProjectFromDB(sqlcProject)), nil
}

func (repo *projectRepo) GetByIDInternal(ctx context.Context, projectID uuid.UUID) (*models.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "GetByIDInternal")
	span.SetAttributes(attribute.String("project.id", projectID.String()))
	defer span.End()
	sqlcProject, err := database.Queries(ctx, repo.q).GetProjectByIDInternal(ctx, projectID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	span.SetAttributes(attribute.String("project.name", sqlcProject.ProjectName))
	return new(mapProjectFromDB(sqlcProject)), nil
}

func (repo *projectRepo) IsOwnerOf(ctx context.Context, projectID, ownerID uuid.UUID) (bool, error) {
	ctx, span := repo.tracer.Start(ctx, "IsOwnerOf")
	span.SetAttributes(attribute.String("project.owner_id", ownerID.String()))
	span.SetAttributes(attribute.String("project.id", projectID.String()))
	defer span.End()
	isOwner, err := database.Queries(ctx, repo.q).IsOwnerOf(ctx, sqlc.IsOwnerOfParams{
		OwnerID: ownerID,
		ID:      projectID,
	})
	if err != nil {
		return false, repo.dbe(err)
	}
	return isOwner, nil
}

func (repo *projectRepo) List(ctx context.Context, ownerID uuid.UUID) ([]models.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "List")
	span.SetAttributes(attribute.String("project.owner_id", ownerID.String()))
	defer span.End()
	sqlcProjects, err := database.Queries(ctx, repo.q).ListProjects(ctx, ownerID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	span.SetAttributes(attribute.Int("project.count", len(sqlcProjects)))
	return xslices.MapSlice(sqlcProjects, mapProjectFromDB), nil
}

func (repo *projectRepo) Update(ctx context.Context, toUpdate models.Project, ownerID uuid.UUID) (*models.Project, error) {
	ctx, span := repo.tracer.Start(ctx, "Update")
	span.SetAttributes(attribute.String("project.owner_id", ownerID.String()))
	span.SetAttributes(attribute.String("project.id", toUpdate.ID.String()))
	defer span.End()
	sqlcProject, err := database.Queries(ctx, repo.q).UpdateProject(ctx, sqlc.UpdateProjectParams{
		ID:          toUpdate.ID,
		OwnerID:     ownerID,
		ProjectName: toUpdate.ProjectName,
		Domain:      toUpdate.Domain,
		Metadata:    toUpdate.Metadata,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapProjectFromDB(sqlcProject)), nil
}

func (repo *projectRepo) Delete(ctx context.Context, projectID, ownerID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "Delete")
	span.SetAttributes(attribute.String("project.owner_id", ownerID.String()))
	span.SetAttributes(attribute.String("project.id", projectID.String()))
	defer span.End()
	affectedRows, err := database.Queries(ctx, repo.q).DeleteProject(ctx, sqlc.DeleteProjectParams{
		ID:      projectID,
		OwnerID: ownerID,
	})
	if err != nil {
		return repo.dbe(err)
	}
	if affectedRows == 0 {
		return fun.ErrNotFound("project not found")
	}
	return nil
}
