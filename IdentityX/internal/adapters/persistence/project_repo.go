package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/project"
	"GoAuth/internal/ports/outbound"
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type projectRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
}

var _ outbound.ProjectRepository = (*projectRepo)(nil)

func NewProjectRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) outbound.ProjectRepository {
	return &projectRepo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func mapUpdateProjectsRowFromDB(dst *project.Project, src *sqlc.UpdateProjectRow) {
	dst.ID = src.ID
	dst.ProjectName = src.ProjectName
	dst.OwnerID = src.OwnerID
	dst.Metadata = src.Metadata
	dst.IsActive = src.IsActive
	dst.PubKey = src.PubKey
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
}

func mapListProjectsRowFromDB(dst *project.Project, src *sqlc.ListProjectsRow) {
	dst.ID = src.ID
	dst.ProjectName = src.ProjectName
	dst.OwnerID = src.OwnerID
	dst.Metadata = src.Metadata
	dst.IsActive = src.IsActive
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
}

func mapGetProjectByIDRowFromDB(dst *project.Project, src *sqlc.GetProjectByIdRow) {
	dst.ID = src.ID
	dst.ProjectName = src.ProjectName
	dst.OwnerID = src.OwnerID
	dst.Metadata = src.Metadata
	dst.IsActive = src.IsActive
	dst.PubKey = src.PubKey
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
}

func mapProjectRowFromDB(dst *project.Project, src *sqlc.CreateProjectRow) {
	dst.ID = src.ID
	dst.ProjectName = src.ProjectName
	dst.OwnerID = src.OwnerID
	dst.Metadata = src.Metadata
	dst.IsActive = src.IsActive
	dst.PubKey = src.PubKey
	dst.PrivKey = []byte(src.PrivKey)
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
}

func (r projectRepo) Create(ctx context.Context, project project.Project) (*project.Project, error) {
	ctx, span := r.tracer.Start(ctx, "ProjectRepo.Create",
		trace.WithAttributes(
			attribute.String("project.owner_id", project.OwnerID.String()),
			attribute.String("project.name", project.ProjectName),
		),
	)
	defer span.End()

	sqlcProject, err := r.q.CreateProject(ctx, sqlc.CreateProjectParams{
		ProjectName: project.ProjectName,
		OwnerID:     project.OwnerID,
		Metadata:    project.Metadata,
		IsActive:    project.IsActive,
		PubKey:      project.PubKey,
		PrivKey:     string(project.PrivKey),
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(attribute.String("project.id", sqlcProject.ID.String()))

	mapProjectRowFromDB(&project, &sqlcProject)
	return &project, nil
}

func (r projectRepo) GetByID(ctx context.Context, projectID, ownerID uuid.UUID) (*project.Project, error) {
	ctx, span := r.tracer.Start(ctx, "ProjectRepo.GetByID",
		trace.WithAttributes(
			attribute.String("project.owner_id", ownerID.String()),
			attribute.String("project.id", projectID.String()),
		),
	)
	defer span.End()

	sqlcProject, err := r.q.GetProjectById(ctx, sqlc.GetProjectByIdParams{
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
	mapGetProjectByIDRowFromDB(&proj, &sqlcProject)
	return &proj, nil
}

func (r projectRepo) GetPublicKeyByID(ctx context.Context, projectID uuid.UUID) (string, error) {
	ctx, span := r.tracer.Start(ctx, "ProjectRepo.GetPublicKeyByID",
		trace.WithAttributes(
			attribute.String("project.id", projectID.String()),
		),
	)
	defer span.End()

	pub, err := r.q.GetProjectPublicKeyById(ctx, projectID)
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return "", sqlcErr
	}

	return pub, nil
}

func (r projectRepo) List(ctx context.Context, ownerID uuid.UUID) ([]project.Project, error) {
	ctx, span := r.tracer.Start(ctx, "ProjectRepo.List",
		trace.WithAttributes(
			attribute.String("project.owner_id", ownerID.String()),
		),
	)
	defer span.End()

	sqlcProjects, err := r.q.ListProjects(ctx, ownerID)
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(attribute.Int("project.count", len(sqlcProjects)))

	projects := make([]project.Project, 0, len(sqlcProjects))
	for _, sqlcProject := range sqlcProjects {
		var proj project.Project
		mapListProjectsRowFromDB(&proj, &sqlcProject)
		projects = append(projects, proj)
	}

	return projects, nil
}

func (r projectRepo) Update(ctx context.Context, project project.Project, ownerID uuid.UUID) (*project.Project, error) {
	ctx, span := r.tracer.Start(ctx, "ProjectRepo.Update",
		trace.WithAttributes(
			attribute.String("project.owner_id", ownerID.String()),
			attribute.String("project.id", project.ID.String()),
		),
	)
	defer span.End()

	sqlcProject, err := r.q.UpdateProject(ctx, sqlc.UpdateProjectParams{
		ID:          project.ID,
		OwnerID:     ownerID,
		ProjectName: project.ProjectName,
		Metadata:    project.Metadata,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	mapUpdateProjectsRowFromDB(&project, &sqlcProject)
	return &project, nil
}

func (r projectRepo) Delete(ctx context.Context, projectID, ownerID uuid.UUID) error {
	ctx, span := r.tracer.Start(ctx, "ProjectRepo.Delete",
		trace.WithAttributes(
			attribute.String("project.owner_id", ownerID.String()),
			attribute.String("project.id", projectID.String()),
		),
	)
	defer span.End()

	err := r.q.DeleteProject(ctx, sqlc.DeleteProjectParams{
		ID:      projectID,
		OwnerID: ownerID,
	})

	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}

	return nil
}
