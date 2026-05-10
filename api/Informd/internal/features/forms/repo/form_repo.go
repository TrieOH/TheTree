package repo

import (
	"Informd/contracts"
	"Informd/internal/platform/database/sqlc"
	"Informd/internal/shared/ports"
	"context"
	"lib/database"
	"lib/errx"
	"lib/xslices"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type formRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    *errx.DBHandler
}

var _ ports.FormsRepo = (*formRepo)(nil)

func NewFormRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer, dbe *errx.DBHandler) ports.FormsRepo {
	return &formRepo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    dbe,
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
		return nil, repo.dbe.DB(err, "form")
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
		return nil, repo.dbe.DB(err, "form")
	}
	forms := xslices.MapSlice(sqlcForms, mapForm)
	forms, err = contracts.FilterForms(forms, params)
	if err != nil {
		return nil, err
	}
	contracts.SortForms(forms, params)
	return forms, nil
}
