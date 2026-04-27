package forms

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

type formRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ ports.FormsRepo = (*formRepo)(nil)

func NewFormRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.FormsRepo {
	return &formRepo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func (repo *formRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func (repo *formRepo) span(ctx context.Context, op string) (context.Context, trace.Span) {
	return repo.tracer.Start(ctx, "FormsRepo."+op)
}

func mapForm(src sqlc.Form) contracts.Form {
	return contracts.Form{
		ID:          src.ID,
		NamespaceID: src.NamespaceID,
		OwnerID:     src.OwnerID,
		Name:        src.Name,
		Status:      contracts.FormStatus(src.Status),
		OpenedAt:    src.OpenedAt,
		ClosedAt:    src.ClosedAt,
		ArchivedAt:  src.ArchivedAt,
		CreatedAt:   src.CreatedAt,
		UpdatedAt:   src.UpdatedAt,
	}
}

func (repo *formRepo) Create(ctx context.Context, toCreate contracts.Form) (*contracts.Form, error) {
	ctx, span := repo.span(ctx, "Create")
	defer span.End()
	sqlcForm, err := repo.queries(ctx).CreateForm(ctx, sqlc.CreateFormParams{
		NamespaceID: toCreate.NamespaceID,
		OwnerID:     toCreate.OwnerID,
		Name:        toCreate.Name,
		Status:      string(toCreate.Status),
	})
	if err != nil {
		return nil, errx.DB(err, "form")
	}
	return new(mapForm(sqlcForm)), nil
}

func (repo *formRepo) GetByID(ctx context.Context, id uuid.UUID) (*contracts.Form, error) {
	ctx, span := repo.span(ctx, "GetByID")
	defer span.End()
	sqlcForm, err := repo.queries(ctx).GetFormByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return new(mapForm(sqlcForm)), nil
}

func (repo *formRepo) BulkGet(ctx context.Context, ids []uuid.UUID) ([]contracts.Form, error) {
	ctx, span := repo.span(ctx, "List")
	defer span.End()
	sqlcForm, err := repo.queries(ctx).BulkGetForms(ctx, ids)
	if err != nil {
		return nil, errx.DB(err, "form")
	}
	return xslices.MapSlice(sqlcForm, mapForm), nil
}
