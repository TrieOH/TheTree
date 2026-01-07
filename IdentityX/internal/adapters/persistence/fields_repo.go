package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/field"
	"GoAuth/internal/ports/outbound"
	"context"
	"database/sql"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type schemaFieldsRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
}

func (repo *schemaFieldsRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(txKeyValue).(*sql.Tx); ok {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

var _ outbound.SchemaFieldsRepository = (*schemaFieldsRepo)(nil)

func NewFieldsRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) outbound.SchemaFieldsRepository {
	return &schemaFieldsRepo{
		q:      q,
		log:    l,
		tracer: tracer,
	}
}

func mapSchemaFieldFromDB(dst *field.Field, src *sqlc.SchemaField) {
	dst.ObjectID = src.ObjectID
	dst.ID = src.ID
	dst.SchemaID = src.SchemaID
	dst.SchemaVersionID = src.SchemaVersionID
	dst.Key = src.Key
	dst.Type = field.Type(src.Type)
	dst.Owner = field.Owner(src.Owner)
	dst.Title = src.Title
	dst.Description = src.Description
	dst.Placeholder = src.Placeholder
	dst.Required = src.Required
	dst.Mutable = src.Mutable
	dst.DefaultValue = src.DefaultValue
	dst.Position = src.Position
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
}

func (repo *schemaFieldsRepo) Create(ctx context.Context, toCreate field.Field) (*field.Field, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.Create",
		trace.WithAttributes(
			attribute.String("field.schema_id", toCreate.SchemaID.String()),
			attribute.String("field.schema_version_id", toCreate.SchemaVersionID.String()),
		),
	)
	defer span.End()

	sqlcSchemaField, err := repo.queries(ctx).CreateSchemaField(ctx, sqlc.CreateSchemaFieldParams{
		Key:             toCreate.Key,
		Type:            sqlc.FieldType(toCreate.Type),
		Owner:           sqlc.FieldOwner(toCreate.Owner),
		Title:           toCreate.Title,
		Description:     toCreate.Description,
		Placeholder:     toCreate.Placeholder,
		Required:        toCreate.Required,
		Mutable:         toCreate.Mutable,
		DefaultValue:    toCreate.DefaultValue,
		Position:        toCreate.Position,
		SchemaVersionID: toCreate.SchemaVersionID,
		SchemaID:        toCreate.SchemaID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(
		attribute.String("field.id", sqlcSchemaField.ID.String()),
		attribute.String("field.object_id", sqlcSchemaField.ObjectID.String()),
	)

	var newSchemaField field.Field
	mapSchemaFieldFromDB(&newSchemaField, &sqlcSchemaField)
	return &newSchemaField, nil
}

func (repo *schemaFieldsRepo) GetByVersionID(ctx context.Context, schemaVersionID uuid.UUID) ([]field.Field, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaFieldsRepo.GetByVersionID",
		trace.WithAttributes(
			attribute.String("field.schema_version_id", schemaVersionID.String()),
		),
	)
	defer span.End()

	sqlcFields, err := repo.queries(ctx).GetFieldsByVersionID(ctx, schemaVersionID)
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(
		attribute.Int("count", len(sqlcFields)),
	)

	fields := make([]field.Field, 0, len(sqlcFields))
	for _, sqlcField := range sqlcFields {
		var newSchemaField field.Field
		mapSchemaFieldFromDB(&newSchemaField, &sqlcField)
		fields = append(fields, newSchemaField)
	}

	return fields, nil
}

func (repo *schemaFieldsRepo) Update(ctx context.Context, toUpdate field.Field) error {
	// TODO Implement me!
	return apierr.ErrInternal.WithMsg("functionality not implemented").WithID(apierr.SystemUnimplemented)
}

func (repo *schemaFieldsRepo) SetOptions(ctx context.Context, options []field.Option) error {
	// TODO Implement me!
	return apierr.ErrInternal.WithMsg("functionality not implemented").WithID(apierr.SystemUnimplemented)
}

func (repo *schemaFieldsRepo) SetRequiredRules(ctx context.Context, rules []field.RequiredRule) error {
	// TODO Implement me!
	return apierr.ErrInternal.WithMsg("functionality not implemented").WithID(apierr.SystemUnimplemented)
}

func (repo *schemaFieldsRepo) SetVisibilityRules(ctx context.Context, rules []field.VisibilityRule) error {
	// TODO Implement me!
	return apierr.ErrInternal.WithMsg("functionality not implemented").WithID(apierr.SystemUnimplemented)
}

func (repo *schemaFieldsRepo) Delete(ctx context.Context, fieldID uuid.UUID) error {
	// TODO Implement me!
	return apierr.ErrInternal.WithMsg("functionality not implemented").WithID(apierr.SystemUnimplemented)
}
