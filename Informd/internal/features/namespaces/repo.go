package namespaces

import (
	"Informd/internal/platform/database"
	"Informd/internal/platform/database/sqlc"
	"Informd/internal/shared/contracts"
	"Informd/internal/shared/errx"
	"Informd/internal/shared/ports"
	"Informd/internal/shared/xslices"
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

var _ ports.NamespaceRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.NamespaceRepo {
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

func (repo *repo) span(ctx context.Context, op string) (context.Context, trace.Span) {
	return repo.tracer.Start(ctx, "ProjectsRepo."+op)
}

func mapProject(src sqlc.Namespace) contracts.Namespace {
	return contracts.Namespace{
		ID:        src.ID,
		OwnerID:   src.OwnerID,
		Name:      src.Name,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
	}
}

func (repo *repo) Create(ctx context.Context, toCreate contracts.Namespace) (*contracts.Namespace, error) {
	ctx, span := repo.span(ctx, "Create")
	defer span.End()

	sqlcProject, err := repo.queries(ctx).CreateNamespace(ctx, sqlc.CreateNamespaceParams{
		OwnerID: toCreate.OwnerID,
		Name:    toCreate.Name,
	})
	if err != nil {
		return nil, errx.DB(err, "project")
	}
	return new(mapProject(sqlcProject)), nil
}

func (repo *repo) GetByID(ctx context.Context, id uuid.UUID) (*contracts.Namespace, error) {
	ctx, span := repo.span(ctx, "GetByID")
	defer span.End()

	sqlcProject, err := repo.queries(ctx).GetNamespaceByID(ctx, id)
	if err != nil {
		return nil, errx.DB(err, "project")
	}
	return new(mapProject(sqlcProject)), nil
}
func (repo *repo) GetByName(ctx context.Context, name string, ownerID uuid.UUID) (*contracts.Namespace, error) {
	ctx, span := repo.span(ctx, "GetByName")
	defer span.End()

	sqlcProject, err := repo.queries(ctx).GetNamespaceByName(ctx, sqlc.GetNamespaceByNameParams{
		OwnerID: ownerID,
		Name:    name,
	})
	if err != nil {
		return nil, errx.DB(err, "project")
	}
	return new(mapProject(sqlcProject)), nil
}

func (repo *repo) List(ctx context.Context, ownerID uuid.UUID) ([]contracts.Namespace, error) {
	ctx, span := repo.span(ctx, "List")
	defer span.End()
	sqlcProjects, err := repo.queries(ctx).ListNamespaceByOwner(ctx, ownerID)
	if err != nil {
		return nil, errx.DB(err, "project")
	}
	return xslices.MapSlice(sqlcProjects, mapProject), nil
}

func (repo *repo) ListByIDs(ctx context.Context, ids []string) ([]contracts.Namespace, error) {
	ctx, span := repo.span(ctx, "ListByIDs")
	defer span.End()
	uuids := xslices.FilterMap(ids, uuid.Parse)
	sqlcProjects, err := repo.queries(ctx).ListNamespaceByIDs(ctx, uuids)
	if err != nil {
		return nil, errx.DB(err, "project")
	}
	return xslices.MapSlice(sqlcProjects, mapProject), nil
}
