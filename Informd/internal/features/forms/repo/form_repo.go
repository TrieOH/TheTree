package repo

import (
	"Informd/internal/platform/database"
	"Informd/internal/platform/database/sqlc"
	"Informd/internal/shared/contracts"
	"Informd/internal/shared/errx"
	"Informd/internal/shared/ports"
	"Informd/internal/shared/xslices"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type formRepo struct {
	q      *sqlc.Queries
	db     *pgxpool.Pool
	log    *zap.Logger
	tracer trace.Tracer
}

var _ ports.FormsRepo = (*formRepo)(nil)

func NewFormRepo(q *sqlc.Queries, db *pgxpool.Pool, log *zap.Logger, tracer trace.Tracer) ports.FormsRepo {
	return &formRepo{
		q:      q,
		db:     db,
		log:    log,
		tracer: tracer,
	}
}

func mapForm(src sqlc.Form) contracts.Form {
	return contracts.Form{
		ID:          src.ID,
		NamespaceID: src.NamespaceID,
		OwnerID:     src.OwnerID,
		Title:       src.Name,
		Status:      contracts.FormStatus(src.Status),
		OpenedAt:    src.OpenedAt,
		ClosedAt:    src.ClosedAt,
		ArchivedAt:  src.ArchivedAt,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
	}
}

func (repo *formRepo) Create(ctx context.Context, toCreate contracts.Form) (*contracts.Form, error) {
	ctx, span := database.Span(ctx, repo.tracer, "Create")
	defer span.End()
	sqlcForm, err := database.Queries(ctx, repo.q).CreateForm(ctx, sqlc.CreateFormParams{
		NamespaceID: toCreate.NamespaceID,
		OwnerID:     toCreate.OwnerID,
		Name:        toCreate.Title,
		Status:      string(toCreate.Status),
	})
	if err != nil {
		return nil, errx.DB(err, "form")
	}
	return new(mapForm(sqlcForm)), nil
}

func (repo *formRepo) GetByID(ctx context.Context, id uuid.UUID) (*contracts.Form, error) {
	ctx, span := database.Span(ctx, repo.tracer, "GetByID")
	defer span.End()
	sqlcForm, err := database.Queries(ctx, repo.q).GetFormByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return new(mapForm(sqlcForm)), nil
}

func (repo *formRepo) BulkGet(ctx context.Context, ids []uuid.UUID, params contracts.BulkGetParams) ([]contracts.Form, error) {
	ctx, span := database.Span(ctx, repo.tracer, "BulkGet")
	defer span.End()
	sqlcForms, err := database.Queries(ctx, repo.q).BulkGetForms(ctx, ids)
	if err != nil {
		return nil, errx.DB(err, "form")
	}
	forms := xslices.MapSlice(sqlcForms, mapForm)
	forms, err = contracts.FilterForms(forms, params)
	if err != nil {
		return nil, err
	}
	contracts.SortForms(forms, params)
	return forms, nil
}
