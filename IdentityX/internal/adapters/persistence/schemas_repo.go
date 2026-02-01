package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/ports/outbounds"
	"context"

	"github.com/MintzyG/fail"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type schemaRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
}

func (repo *schemaRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(transactions.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

var _ outbounds.SchemaRepository = (*schemaRepo)(nil)

func NewSchemaRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) outbounds.SchemaRepository {
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

func (repo *schemaRepo) Draft(ctx context.Context, toDraft schema.Schema) (*schema.Schema, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaRepo.Draft",
		trace.WithAttributes(
			attribute.String("schema.project_id", toDraft.ProjectID.String()),
			attribute.String("schema.type", string(toDraft.Type)),
		),
	)
	defer span.End()

	sqlcSchema, err := repo.queries(ctx).DraftSchema(ctx, sqlc.DraftSchemaParams{
		ProjectID: toDraft.ProjectID,
		Title:     toDraft.Title,
		FlowID:    toDraft.FlowID,
		Type:      sqlc.SchemaType(toDraft.Type),
	})
	if err != nil {
		return nil, fail.From(err)
	}

	span.SetAttributes(attribute.String("schema.id", sqlcSchema.ID.String()))

	var createdSchema schema.Schema
	mapSchemaFromDB(&createdSchema, &sqlcSchema)
	return &createdSchema, nil
}

func (repo *schemaRepo) Publish(ctx context.Context, toPublish schema.Schema) error {
	ctx, span := repo.tracer.Start(ctx, "SchemaRepo.Publish",
		trace.WithAttributes(
			attribute.String("schema.id", toPublish.ID.String()),
			attribute.String("schema.project_id", toPublish.ProjectID.String()),
		),
	)
	defer span.End()

	if err := repo.queries(ctx).PublishSchema(ctx, sqlc.PublishSchemaParams{
		ID:        toPublish.ID,
		ProjectID: toPublish.ProjectID,
	}); err != nil {
		return fail.From(err)
	}

	return nil
}

func (repo *schemaRepo) Archive(ctx context.Context, toArchive schema.Schema) error {
	ctx, span := repo.tracer.Start(ctx, "SchemaRepo.Archive",
		trace.WithAttributes(
			attribute.String("schema.id", toArchive.ID.String()),
			attribute.String("schema.project_id", toArchive.ProjectID.String()),
		),
	)
	defer span.End()

	if err := repo.queries(ctx).ArchiveSchema(ctx, sqlc.ArchiveSchemaParams{
		ID:        toArchive.ID,
		ProjectID: toArchive.ProjectID,
	}); err != nil {
		return fail.From(err)
	}

	return nil
}

func (repo *schemaRepo) Delete(ctx context.Context, toDelete schema.Schema) error {
	ctx, span := repo.tracer.Start(ctx, "SchemaRepo.Delete",
		trace.WithAttributes(
			attribute.String("schema.id", toDelete.ID.String()),
			attribute.String("schema.project_id", toDelete.ProjectID.String()),
		),
	)
	defer span.End()

	if err := repo.queries(ctx).DeleteSchema(ctx, sqlc.DeleteSchemaParams{
		ID:        toDelete.ID,
		ProjectID: toDelete.ProjectID,
	}); err != nil {
		return fail.From(err)
	}

	return nil
}

func (repo *schemaRepo) Exists(ctx context.Context, toCheck schema.Schema) (bool, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaRepo.Exists",
		trace.WithAttributes(
			attribute.String("schema.project_id", toCheck.ProjectID.String()),
		),
	)
	defer span.End()

	exists, err := repo.queries(ctx).SchemaExists(ctx, sqlc.SchemaExistsParams{
		ProjectID: toCheck.ProjectID,
		FlowID:    toCheck.FlowID,
		Type:      sqlc.SchemaType(toCheck.Type),
	})
	if err != nil {
		return false, fail.From(err)
	}

	return exists, nil
}

func (repo *schemaRepo) BelongsToProject(ctx context.Context, toCheck schema.Schema) (bool, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaRepo.BelongsToProject",
		trace.WithAttributes(
			attribute.String("schema.project_id", toCheck.ProjectID.String()),
			attribute.String("schema.id", toCheck.ID.String()),
		),
	)
	defer span.End()

	belongs, err := repo.queries(ctx).SchemaBelongsToProject(ctx, sqlc.SchemaBelongsToProjectParams{
		ID:        toCheck.ID,
		ProjectID: toCheck.ProjectID,
	})
	if err != nil {
		return false, fail.From(err)
	}

	return belongs, nil
}

func (repo *schemaRepo) FindByID(ctx context.Context, schemaID uuid.UUID, projectID uuid.UUID) (*schema.Schema, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaRepo.FindByID",
		trace.WithAttributes(
			attribute.String("schema.id", schemaID.String()),
			attribute.String("schema.project_id", projectID.String()),
		),
	)
	defer span.End()

	sqlcSchema, err := repo.queries(ctx).GetSchema(ctx, sqlc.GetSchemaParams{
		ID:        schemaID,
		ProjectID: projectID,
	})
	if err != nil {
		return nil, fail.From(err)
	}

	span.SetAttributes(attribute.String("schema.type", string(sqlcSchema.Type)))

	var foundSchema schema.Schema
	mapSchemaFromDB(&foundSchema, &sqlcSchema)
	return &foundSchema, nil
}

func (repo *schemaRepo) FindByFlowIDAndType(ctx context.Context, flowID string, schemaType schema.Type, projectID uuid.UUID) (*schema.Schema, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaRepo.FindByFlowIDAndType",
		trace.WithAttributes(
			attribute.String("schema.project_id", projectID.String()),
		),
	)
	defer span.End()

	sqlcSchema, err := repo.queries(ctx).GetSchemaByFlowIDAndType(ctx, sqlc.GetSchemaByFlowIDAndTypeParams{
		FlowID:    flowID,
		Type:      sqlc.SchemaType(schemaType),
		ProjectID: projectID,
	})
	if err != nil {
		return nil, fail.From(err).WithArgs("schema")
	}

	span.SetAttributes(attribute.String("schema.id", sqlcSchema.ID.String()))

	var foundSchema schema.Schema
	mapSchemaFromDB(&foundSchema, &sqlcSchema)
	return &foundSchema, nil
}

func (repo *schemaRepo) List(ctx context.Context, projectID uuid.UUID) ([]schema.Schema, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaRepo.List",
		trace.WithAttributes(
			attribute.String("schema.project_id", projectID.String()),
		),
	)
	defer span.End()

	sqlcSchemas, err := repo.queries(ctx).ListSchemas(ctx, projectID)
	if err != nil {
		return nil, fail.From(err)
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

func (repo *schemaRepo) SetVersion(ctx context.Context, toUpdate schema.Schema) error {
	ctx, span := repo.tracer.Start(ctx, "SchemaRepo.SetVersion",
		trace.WithAttributes(
			attribute.String("schema.project_id", toUpdate.ProjectID.String()),
			attribute.String("schema.id", toUpdate.ID.String()),
		),
	)
	if toUpdate.CurrentVersionID != nil {
		span.SetAttributes(attribute.String("schema.current_version_id", toUpdate.CurrentVersionID.String()))
	}
	defer span.End()

	err := repo.queries(ctx).SetSchemaVersion(ctx, sqlc.SetSchemaVersionParams{
		CurrentVersionID: toUpdate.CurrentVersionID,
		ID:               toUpdate.ID,
		ProjectID:        toUpdate.ProjectID,
	})
	if err != nil {
		return fail.From(err)
	}

	return nil
}

func (repo *schemaRepo) GetIDsFromProjectID(ctx context.Context, projectID uuid.UUID) ([]uuid.UUID, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaRepo.GetIDsFromProjectID",
		trace.WithAttributes(attribute.String("project_id", projectID.String())),
	)
	defer span.End()

	IDs, err := repo.queries(ctx).GetSchemaIDsFromProjectID(ctx, projectID)
	if err != nil {
		return nil, fail.From(err)
	}

	span.SetAttributes(attribute.Int("schema.count", len(IDs)))

	return IDs, nil
}
