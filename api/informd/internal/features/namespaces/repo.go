package namespaces

import (
	"Informd/internal/database/sqlc"
	"Informd/models"
	"Informd/ports"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type repo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.NamespaceRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.NamespaceRepo {
	return &repo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("namespace"),
	}
}

func mapNamespace(src sqlc.Namespace) models.Namespace {
	return models.Namespace{
		ID:        src.ID,
		OwnerID:   src.OwnerID,
		Name:      src.Name,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
	}
}

func (repo *repo) Create(ctx context.Context, toCreate models.Namespace) (*models.Namespace, error) {
	ctx, span := repo.tracer.Start(ctx, "Create")
	defer span.End()
	sqlcProject, err := database.Queries(ctx, repo.q).CreateNamespace(ctx, sqlc.CreateNamespaceParams{
		OwnerID: toCreate.OwnerID,
		Name:    toCreate.Name,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapNamespace(sqlcProject)), nil
}

func (repo *repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Namespace, error) {
	ctx, span := repo.tracer.Start(ctx, "GetByID")
	defer span.End()
	sqlcProject, err := database.Queries(ctx, repo.q).GetNamespaceByID(ctx, id)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapNamespace(sqlcProject)), nil
}
func (repo *repo) GetByName(ctx context.Context, name string, ownerID uuid.UUID) (*models.Namespace, error) {
	ctx, span := repo.tracer.Start(ctx, "GetByName")
	defer span.End()
	sqlcProject, err := database.Queries(ctx, repo.q).GetNamespaceByName(ctx, sqlc.GetNamespaceByNameParams{
		OwnerID: ownerID,
		Name:    name,
	})
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapNamespace(sqlcProject)), nil
}

func (repo *repo) BulkGet(ctx context.Context, ids []uuid.UUID) ([]models.Namespace, error) {
	ctx, span := repo.tracer.Start(ctx, "BulkGet")
	defer span.End()
	sqlcForm, err := database.Queries(ctx, repo.q).BulkGetNamespaces(ctx, ids)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(sqlcForm, mapNamespace), nil
}
