package repo

import (
	"GoAuth/internal/apierr"
	"GoAuth/internal/models"
	"GoAuth/internal/sqlc"
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type ProjectRepo interface {
	Create(ctx context.Context, project models.Project) (*models.Project, error)

	GetByID(ctx context.Context, projectID, ownerID uuid.UUID) (*models.Project, error)
	GetKeysByID(ctx context.Context, projectID, ownerID uuid.UUID) (*models.ProjectKeys, error)
	GetPublicKeyByID(ctx context.Context, projectID uuid.UUID) (string, error)

	List(ctx context.Context, ownerID uuid.UUID) ([]models.Project, error)

	Update(ctx context.Context, project models.Project, ownerID uuid.UUID) (*models.Project, error)
	Delete(ctx context.Context, projectID, ownerID uuid.UUID) error

	AdminGetByID(ctx context.Context, projectID uuid.UUID) (*models.Project, error)
}

type projectRepo struct {
	q   *sqlc.Queries
	log *zap.Logger
}

func NewProjectRepo(q *sqlc.Queries, log *zap.Logger) ProjectRepo {
	return &projectRepo{
		q:   q,
		log: log,
	}
}

func mapAdminGetProjectByIDRowFromDB(dst *models.Project, src *sqlc.AdminGetProjectByIdRow) {
	dst.ID = src.ID
	dst.ProjectName = src.ProjectName
	dst.OwnerID = src.OwnerID
	dst.Metadata = src.Metadata
	dst.IsActive = src.IsActive
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
}

func mapUpdateProjectsRowFromDB(dst *models.Project, src *sqlc.UpdateProjectRow) {
	dst.ID = src.ID
	dst.ProjectName = src.ProjectName
	dst.OwnerID = src.OwnerID
	dst.Metadata = src.Metadata
	dst.IsActive = src.IsActive
	dst.PubKey = src.PubKey
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
}

func mapListProjectsRowFromDB(dst *models.Project, src *sqlc.ListProjectsRow) {
	dst.ID = src.ID
	dst.ProjectName = src.ProjectName
	dst.OwnerID = src.OwnerID
	dst.Metadata = src.Metadata
	dst.IsActive = src.IsActive
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
}

func mapGetProjectByIDRowFromDB(dst *models.Project, src *sqlc.GetProjectByIdRow) {
	dst.ID = src.ID
	dst.ProjectName = src.ProjectName
	dst.OwnerID = src.OwnerID
	dst.Metadata = src.Metadata
	dst.IsActive = src.IsActive
	dst.PubKey = src.PubKey
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
}

func mapProjectRowFromDB(dst *models.Project, src *sqlc.CreateProjectRow) {
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

func (r projectRepo) Create(ctx context.Context, project models.Project) (*models.Project, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "ProjectRepo.Create",
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

func (r projectRepo) GetByID(ctx context.Context, projectID, ownerID uuid.UUID) (*models.Project, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "ProjectRepo.GetByID",
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

	var project models.Project
	mapGetProjectByIDRowFromDB(&project, &sqlcProject)
	return &project, nil
}

func (r projectRepo) GetKeysByID(ctx context.Context, projectID, ownerID uuid.UUID) (*models.ProjectKeys, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "ProjectRepo.GetKeysByID",
		trace.WithAttributes(
			attribute.String("project.owner_id", ownerID.String()),
			attribute.String("project.id", projectID.String()),
		),
	)
	defer span.End()

	keys, err := r.q.GetProjectKeysById(ctx, sqlc.GetProjectKeysByIdParams{
		ID:      projectID,
		OwnerID: ownerID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	return &models.ProjectKeys{
		PubKey:  keys.PubKey,
		PrivKey: []byte(keys.PrivKey),
	}, nil
}

func (r projectRepo) GetPublicKeyByID(ctx context.Context, projectID uuid.UUID) (string, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "ProjectRepo.GetPublicKeyByID",
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

func (r projectRepo) List(ctx context.Context, ownerID uuid.UUID) ([]models.Project, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "ProjectRepo.List",
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

	projects := make([]models.Project, 0, len(sqlcProjects))
	for _, sqlcProject := range sqlcProjects {
		var project models.Project
		mapListProjectsRowFromDB(&project, &sqlcProject)
		projects = append(projects, project)
	}

	return projects, nil
}

func (r projectRepo) Update(ctx context.Context, project models.Project, ownerID uuid.UUID) (*models.Project, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "ProjectRepo.Update",
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
	ctx, span := GoAuthRepoTracer.Start(ctx, "ProjectRepo.Delete",
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

func (r projectRepo) AdminGetByID(ctx context.Context, projectID uuid.UUID) (*models.Project, error) {
	ctx, span := GoAuthRepoTracer.Start(ctx, "ProjectRepo.AdminGetByID",
		trace.WithAttributes(
			attribute.String("project.id", projectID.String()),
		),
	)
	defer span.End()

	sqlcProject, err := r.q.AdminGetProjectById(ctx, projectID)
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(
		attribute.String("project.owner_id", sqlcProject.OwnerID.String()),
		attribute.String("project.name", sqlcProject.ProjectName),
	)

	var project models.Project
	mapAdminGetProjectByIDRowFromDB(&project, &sqlcProject)
	return &project, nil
}
