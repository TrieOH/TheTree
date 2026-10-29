package repo

import (
	"Informd/errs"
	"Informd/internal/sqlc"
	"Informd/models"
	"Informd/ports"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.FieldsRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrFieldNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapField(src sqlc.Field) models.Field {
	return models.Field{
		ID:           src.ID,
		StepID:       src.StepID,
		Key:          src.Key,
		Title:        src.Title,
		Description:  src.Description,
		PositionHint: src.PositionHint,
		Required:     src.Required,
		Type:         models.FieldType(src.Type),
		Placeholder:  src.Placeholder,
		DefaultValue: src.DefaultValue,
		Config:       src.Config,
		CreatedAt:    src.CreatedAt,
		UpdatedAt:    src.UpdatedAt,
	}
}

func mapFieldSelectConfig(src sqlc.FieldSelectConfig) models.FieldSelectConfig {
	return models.FieldSelectConfig{
		FieldID:   src.FieldID,
		Behaviour: models.SelectBehaviour(src.Behaviour),
		ValueType: models.SelectValueType(src.ValueType),
		Options:   src.Options,
	}
}

func (repo *Repo) BulkEdit(ctx context.Context, fields []models.Field) error {
	params := xslices.MapSlice(fields, ToBulkEditFieldsParams)
	return database.BatchExec(
		repo.db(ctx).BulkEditFields(ctx, params),
		repo.dbe,
		func(i int) string { return params[i].ID.String() },
	)
}

func ToBulkEditFieldsParams(f models.Field) sqlc.BulkEditFieldsParams {
	return sqlc.BulkEditFieldsParams{
		ID:           f.ID,
		StepID:       f.StepID,
		Key:          f.Key,
		Title:        f.Title,
		Description:  f.Description,
		PositionHint: f.PositionHint,
		Required:     f.Required,
		Type:         string(f.Type),
		Placeholder:  f.Placeholder,
		DefaultValue: f.DefaultValue,
		Config:       f.Config,
	}
}

func (repo *Repo) Create(ctx context.Context, toCreate models.Field) (*models.Field, error) {
	sqlcField, err := repo.db(ctx).CreateField(ctx, sqlc.CreateFieldParams{
		StepID:       toCreate.StepID,
		Key:          toCreate.Key,
		Title:        toCreate.Title,
		Description:  toCreate.Description,
		PositionHint: toCreate.PositionHint,
		Required:     toCreate.Required,
		Type:         string(toCreate.Type),
		Placeholder:  toCreate.Placeholder,
		DefaultValue: toCreate.DefaultValue,
		Config:       toCreate.Config,
	})
	return database.One(sqlcField, err, repo.dbe, mapField)
}

func (repo *Repo) CreateSelectConfig(ctx context.Context, toCreate models.FieldSelectConfig) (*models.FieldSelectConfig, error) {
	sqlcConfig, err := repo.db(ctx).CreateFieldSelectConfig(ctx, sqlc.CreateFieldSelectConfigParams{
		FieldID:   toCreate.FieldID,
		Behaviour: string(toCreate.Behaviour),
		ValueType: string(toCreate.ValueType),
		Options:   toCreate.Options,
	})
	return database.One(sqlcConfig, err, repo.dbe, mapFieldSelectConfig)
}

func (repo *Repo) Delete(ctx context.Context, id uuid.UUID) error {
	err := repo.db(ctx).DeleteField(ctx, id)
	return repo.dbe(err)
}

func (repo *Repo) DeleteSelectConfig(ctx context.Context, fieldID uuid.UUID) error {
	err := repo.db(ctx).DeleteFieldSelectConfig(ctx, fieldID)
	return repo.dbe(err)
}

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.Field, error) {
	sqlcField, err := repo.db(ctx).GetFieldByID(ctx, id)
	return database.One(sqlcField, err, repo.dbe, mapField)
}

func (repo *Repo) GetSelectConfig(ctx context.Context, fieldID uuid.UUID) (*models.FieldSelectConfig, error) {
	sqlcConfig, err := repo.db(ctx).GetFieldSelectConfig(ctx, fieldID)
	return database.One(sqlcConfig, err, repo.dbe, mapFieldSelectConfig)
}

func (repo *Repo) ListByFormID(ctx context.Context, formID uuid.UUID) ([]models.Field, error) {
	rows, err := repo.db(ctx).ListFieldsByFormID(ctx, formID)
	return database.Many(rows, err, repo.dbe, mapField)
}

func (repo *Repo) ListByStepID(ctx context.Context, stepID uuid.UUID) ([]models.Field, error) {
	sqlcFields, err := repo.db(ctx).ListFieldsByStepID(ctx, stepID)
	return database.Many(sqlcFields, err, repo.dbe, mapField)
}

func (repo *Repo) UpdateSelectConfig(ctx context.Context, toUpdate models.FieldSelectConfig) (*models.FieldSelectConfig, error) {
	sqlcConfig, err := repo.db(ctx).UpdateFieldSelectConfig(ctx, sqlc.UpdateFieldSelectConfigParams{
		FieldID:   toUpdate.FieldID,
		Behaviour: string(toUpdate.Behaviour),
		ValueType: string(toUpdate.ValueType),
		Options:   toUpdate.Options,
	})
	return database.One(sqlcConfig, err, repo.dbe, mapFieldSelectConfig)
}
