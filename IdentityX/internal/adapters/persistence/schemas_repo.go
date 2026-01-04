package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/ports/outbound"
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type schemaRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
}

var _ outbound.SchemaRepository = (*schemaRepo)(nil)

func NewSchemaRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) outbound.SchemaRepository {
	return &schemaRepo{
		q:      q,
		log:    l,
		tracer: tracer,
	}
}

func mapSchemaFromDB(dst *schema.Schema, src *sqlc.Schema) {
	dst.ID = src.ID
	dst.ProjectID = src.ProjectID
	dst.Title = src.Title
	dst.FlowID = src.FlowID
	dst.Type = schema.Type(src.Type)
	dst.CurrentVersionID = src.CurrentVersionID
	dst.Status = schema.Status(src.Status)
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
}

func (r schemaRepo) Draft(ctx context.Context, createSchema schema.Schema) (*schema.Schema, error) {
	ctx, span := r.tracer.Start(ctx, "SchemaRepo.Draft",
		trace.WithAttributes(
			attribute.String("schema.project_id", createSchema.ProjectID.String()),
			attribute.String("schema.type", string(createSchema.Type)),
		),
	)
	defer span.End()

	sqlcSchema, err := r.q.DraftSchema(ctx, sqlc.DraftSchemaParams{
		ProjectID: createSchema.ProjectID,
		Title:     createSchema.Title,
		FlowID:    createSchema.FlowID,
		Type:      sqlc.SchemaType(createSchema.Type),
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(attribute.String("schema.id", sqlcSchema.ID.String()))

	var newSchema schema.Schema
	mapSchemaFromDB(&newSchema, &sqlcSchema)
	return &newSchema, nil
}

func (r schemaRepo) Publish(ctx context.Context, schemaID uuid.UUID, projectID uuid.UUID) error {
	ctx, span := r.tracer.Start(ctx, "SchemaRepo.Publish",
		trace.WithAttributes(
			attribute.String("schema.id", schemaID.String()),
			attribute.String("schema.project_id", projectID.String()),
		),
	)
	defer span.End()

	if err := r.q.PublishSchema(ctx, sqlc.PublishSchemaParams{
		ID:        schemaID,
		ProjectID: projectID,
	}); err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return err
	}

	return nil
}

func (r schemaRepo) Archive(ctx context.Context, schemaID uuid.UUID, projectID uuid.UUID) error {
	ctx, span := r.tracer.Start(ctx, "SchemaRepo.Archive",
		trace.WithAttributes(
			attribute.String("schema.id", schemaID.String()),
			attribute.String("schema.project_id", projectID.String()),
		),
	)
	defer span.End()

	if err := r.q.ArchiveSchema(ctx, sqlc.ArchiveSchemaParams{
		ID:        schemaID,
		ProjectID: projectID,
	}); err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return err
	}

	return nil
}

func (r schemaRepo) Delete(ctx context.Context, schemaID uuid.UUID, projectID uuid.UUID) error {
	ctx, span := r.tracer.Start(ctx, "SchemaRepo.Delete",
		trace.WithAttributes(
			attribute.String("schema.id", schemaID.String()),
			attribute.String("schema.project_id", projectID.String()),
		),
	)
	defer span.End()

	if err := r.q.DeleteSchema(ctx, sqlc.DeleteSchemaParams{
		ID:        schemaID,
		ProjectID: projectID,
	}); err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return err
	}

	return nil
}

func (r schemaRepo) FindByID(ctx context.Context, schemaID uuid.UUID, projectID uuid.UUID) (*schema.Schema, error) {
	ctx, span := r.tracer.Start(ctx, "SchemaRepo.FindByID",
		trace.WithAttributes(
			attribute.String("schema.id", schemaID.String()),
			attribute.String("schema.project_id", projectID.String()),
		),
	)
	defer span.End()

	slqcSchema, err := r.q.GetSchema(ctx, sqlc.GetSchemaParams{
		ID:        schemaID,
		ProjectID: projectID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, err
	}

	span.SetAttributes(attribute.String("schema.type", string(slqcSchema.Type)))

	var foundSchema schema.Schema
	mapSchemaFromDB(&foundSchema, &slqcSchema)
	return &foundSchema, nil
}

func (r schemaRepo) List(ctx context.Context, projectID uuid.UUID) ([]schema.Schema, error) {
	ctx, span := r.tracer.Start(ctx, "SchemaRepo.FindByID",
		trace.WithAttributes(
			attribute.String("schema.project_id", projectID.String()),
		),
	)
	defer span.End()

	slqcSchemas, err := r.q.ListSchemas(ctx, projectID)
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, err
	}

	span.SetAttributes(attribute.Int("schema.count", len(slqcSchemas)))

	var schemaList []schema.Schema
	for _, slqcSchema := range slqcSchemas {
		var foundSchema schema.Schema
		mapSchemaFromDB(&foundSchema, &slqcSchema)
		schemaList = append(schemaList, foundSchema)
	}
	return schemaList, nil
}
