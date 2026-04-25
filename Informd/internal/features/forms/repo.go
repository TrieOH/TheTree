package forms

import (
	"Informd/internal/platform/database"
	"Informd/internal/platform/database/sqlc"
	"Informd/internal/shared/contracts"
	"Informd/internal/shared/errx"
	"Informd/internal/shared/ports"
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

var _ ports.FormsRepo = (*repo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.FormsRepo {
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

func mapForm(src *sqlc.Form) *contracts.Form {
	return &contracts.Form{
		ID:        src.ID,
		ProjectID: src.ProjectID,
		OwnerID:   src.OwnerID,
		Title:     src.Title,
		Status:    contracts.FormStatus(src.Status),
		//CurrentVersionID: src.CurrentVersionID,
		CreatedAt:  src.CreatedAt,
		UpdatedAt:  src.UpdatedAt,
		OpenedAt:   src.OpenedAt,
		ClosedAt:   src.ClosedAt,
		ArchivedAt: src.ArchivedAt,
	}
}

func (repo *repo) Create(ctx context.Context, toCreate contracts.Form) (*contracts.Form, error) {
	ctx, span := repo.tracer.Start(ctx, "FormsRepo.Create")
	defer span.End()

	sqlcForm, err := repo.queries(ctx).CreateForm(ctx, sqlc.CreateFormParams{
		ID:        toCreate.ID,
		ProjectID: toCreate.ProjectID,
		OwnerID:   toCreate.OwnerID,
		Title:     toCreate.Title,
		Status:    string(toCreate.Status),
	})
	if err != nil {
		return nil, errx.DB(err, "form")
	}

	return mapForm(&sqlcForm), nil
}

func (repo *repo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]contracts.Form, error) {
	ctx, span := repo.tracer.Start(ctx, "FormsRepo.ListByProject")
	defer span.End()

	sqlcForm, err := repo.queries(ctx).ListFormsByProject(ctx, projectID)
	if err != nil {
		return nil, errx.DB(err, "form")
	}

	out := make([]contracts.Form, 0, len(sqlcForm))
	for _, form := range sqlcForm {
		out = append(out, *mapForm(&form))
	}
	return out, nil
}
