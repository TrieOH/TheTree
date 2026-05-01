package forms

import (
	"Informd/internal/platform/database"
	"Informd/internal/platform/database/sqlc"
	"Informd/internal/shared/contracts"
	"Informd/internal/shared/errx"
	"Informd/internal/shared/ports"
	"Informd/internal/shared/xslices"
	"context"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (repo *formRepo) BulkGet(ctx context.Context, ids []uuid.UUID, params contracts.BulkGetParams) ([]contracts.Form, error) {
	ctx, span := repo.span(ctx, "BulkGet")
	defer span.End()

	isNullOp := params.FilterOp == "is_null" || params.FilterOp == "not_null"

	if params.FilterValue == "" && !isNullOp {
		sqlcForms, err := repo.queries(ctx).BulkGetForms(ctx, ids)
		if err != nil {
			return nil, errx.DB(err, "form")
		}
		return xslices.MapSlice(sqlcForms, mapForm), nil
	}

	query, args, err := buildBulkGetQuery(ids, params)
	if err != nil {
		return nil, err
	}
	rows, err := repo.db.Query(ctx, query, args...)
	if err != nil {
		return nil, errx.DB(err, "form")
	}
	defer rows.Close()

	var forms []contracts.Form
	for rows.Next() {
		var f sqlc.Form
		if err = rows.Scan(
			&f.ID, &f.OwnerID, &f.NamespaceID, &f.Name,
			&f.Status, &f.OpenedAt, &f.ClosedAt, &f.ArchivedAt,
			&f.CreatedAt, &f.UpdatedAt,
		); err != nil {
			return nil, errx.DB(err, "form")
		}
		forms = append(forms, mapForm(f))
	}
	return forms, rows.Err()
}

func buildBulkGetQuery(ids []uuid.UUID, params contracts.BulkGetParams) (string, []any, error) {
	col, ok := contracts.AllowedFilterKeys[params.FilterKey]
	if !ok {
		return "", nil, fmt.Errorf("invalid filter_key: %s", params.FilterKey)
	}
	op, ok := contracts.AllowedOps[params.FilterOp]
	if !ok {
		return "", nil, fmt.Errorf("invalid filter_op: %s", params.FilterOp)
	}

	order := "created_at DESC"
	if strings.EqualFold(params.FilterOrder, "asc") {
		order = "created_at ASC"
	}

	q := sq.StatementBuilder.
		PlaceholderFormat(sq.Dollar).
		Select("id, owner_id, namespace_id, name, status, opened_at, closed_at, archived_at, created_at, updated_at").
		From("forms").
		Where(sq.Expr("id = ANY(?)", ids))

	if params.FilterOp == "is_null" || params.FilterOp == "not_null" {
		q = q.Where(fmt.Sprintf("%s %s", col, op)) // sem placeholder
	} else {
		q = q.Where(fmt.Sprintf("%s %s ?", col, op), params.FilterValue)
	}

	return q.OrderBy(order).ToSql()
}
