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

func (r schemaRepo) Draft(ctx context.Context, newSchema schema.Schema) (*schema.Schema, error) {
	ctx, span := r.tracer.Start(ctx, "SchemaRepo.Draft",
		trace.WithAttributes(
			attribute.String("schema.project_id", newSchema.ProjectID.String()),
			attribute.String("schema.type", string(newSchema.Type)),
		),
	)
	defer span.End()

	sqlcSchema, err := r.q.DraftSchema(ctx, sqlc.DraftSchemaParams{
		ProjectID: newSchema.ProjectID,
		Title:     newSchema.Title,
		FlowID:    newSchema.FlowID,
		Type:      sqlc.SchemaType(newSchema.Type),
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(attribute.String("schema.id", sqlcSchema.ID.String()))

	var createdSchema schema.Schema
	mapSchemaFromDB(&createdSchema, &sqlcSchema)
	return &createdSchema, nil
}

func (r schemaRepo) Publish(ctx context.Context, publishedSchema schema.Schema) error {
	ctx, span := r.tracer.Start(ctx, "SchemaRepo.Publish",
		trace.WithAttributes(
			attribute.String("schema.id", publishedSchema.ID.String()),
			attribute.String("schema.project_id", publishedSchema.ProjectID.String()),
		),
	)
	defer span.End()

	if err := r.q.PublishSchema(ctx, sqlc.PublishSchemaParams{
		ID:        publishedSchema.ID,
		ProjectID: publishedSchema.ProjectID,
	}); err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}

	return nil
}

func (r schemaRepo) Archive(ctx context.Context, archivedSchema schema.Schema) error {
	ctx, span := r.tracer.Start(ctx, "SchemaRepo.Archive",
		trace.WithAttributes(
			attribute.String("schema.id", archivedSchema.ID.String()),
			attribute.String("schema.project_id", archivedSchema.ProjectID.String()),
		),
	)
	defer span.End()

	if err := r.q.ArchiveSchema(ctx, sqlc.ArchiveSchemaParams{
		ID:        archivedSchema.ID,
		ProjectID: archivedSchema.ProjectID,
	}); err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}

	return nil
}

func (r schemaRepo) Delete(ctx context.Context, deletedSchema schema.Schema) error {
	ctx, span := r.tracer.Start(ctx, "SchemaRepo.Delete",
		trace.WithAttributes(
			attribute.String("schema.id", deletedSchema.ID.String()),
			attribute.String("schema.project_id", deletedSchema.ProjectID.String()),
		),
	)
	defer span.End()

	if err := r.q.DeleteSchema(ctx, sqlc.DeleteSchemaParams{
		ID:        deletedSchema.ID,
		ProjectID: deletedSchema.ProjectID,
	}); err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}

	return nil
}

func (r schemaRepo) Exists(ctx context.Context, existsSchema schema.Schema) (bool, error) {
	ctx, span := r.tracer.Start(ctx, "SchemaRepo.Exists",
		trace.WithAttributes(
			attribute.String("schema.project_id", existsSchema.ProjectID.String()),
		),
	)
	defer span.End()

	exists, err := r.q.SchemaExists(ctx, sqlc.SchemaExistsParams{
		ProjectID: existsSchema.ProjectID,
		FlowID:    existsSchema.FlowID,
		Type:      sqlc.SchemaType(existsSchema.Type),
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return false, sqlcErr
	}

	return exists, nil
}

func (r schemaRepo) FindByID(ctx context.Context, schemaID uuid.UUID, projectID uuid.UUID) (*schema.Schema, error) {
	ctx, span := r.tracer.Start(ctx, "SchemaRepo.FindByID",
		trace.WithAttributes(
			attribute.String("schema.id", schemaID.String()),
			attribute.String("schema.project_id", projectID.String()),
		),
	)
	defer span.End()

	sqlcSchema, err := r.q.GetSchema(ctx, sqlc.GetSchemaParams{
		ID:        schemaID,
		ProjectID: projectID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(attribute.String("schema.type", string(sqlcSchema.Type)))

	var foundSchema schema.Schema
	mapSchemaFromDB(&foundSchema, &sqlcSchema)
	return &foundSchema, nil
}

func (r schemaRepo) FindByFlowID(ctx context.Context, flowID string, projectID uuid.UUID) (*schema.Schema, error) {
	ctx, span := r.tracer.Start(ctx, "SchemaRepo.FindByFlowID",
		trace.WithAttributes(
			attribute.String("schema.project_id", projectID.String()),
		),
	)
	defer span.End()

	sqlcSchema, err := r.q.GetSchemaByFlowID(ctx, sqlc.GetSchemaByFlowIDParams{
		FlowID:    flowID,
		ProjectID: projectID,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(attribute.String("schema.type", string(sqlcSchema.Type)))
	span.SetAttributes(attribute.String("schema.id", sqlcSchema.ID.String()))

	var foundSchema schema.Schema
	mapSchemaFromDB(&foundSchema, &sqlcSchema)
	return &foundSchema, nil
}

func (r schemaRepo) List(ctx context.Context, projectID uuid.UUID) ([]schema.Schema, error) {
	ctx, span := r.tracer.Start(ctx, "SchemaRepo.List",
		trace.WithAttributes(
			attribute.String("schema.project_id", projectID.String()),
		),
	)
	defer span.End()

	sqlcSchemas, err := r.q.ListSchemas(ctx, projectID)
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(attribute.Int("schema.count", len(sqlcSchemas)))

	schemaList := make([]schema.Schema, 0, len(sqlcSchemas))
	for _, sqlcSchema := range sqlcSchemas {
		var foundSchema schema.Schema
		mapSchemaFromDB(&foundSchema, &sqlcSchema)
		schemaList = append(schemaList, foundSchema)
	}
	return schemaList, nil
}
