package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/apierr"
	"GoAuth/internal/domain/schema"
	"GoAuth/internal/ports/outbound"
	"context"
	"database/sql"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type schemaVersionRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
}

func (repo *schemaVersionRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(txKeyValue).(*sql.Tx); ok {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

var _ outbound.SchemaVersionRepository = (*schemaVersionRepo)(nil)

func NewSchemaVersionRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) outbound.SchemaVersionRepository {
	return &schemaVersionRepo{
		q:      q,
		log:    l,
		tracer: tracer,
	}
}

func mapSchemaVersionFromDB(dst *schema.Version, src *sqlc.SchemaVersion) {
	dst.ID = src.ID
	dst.SchemaID = src.SchemaID
	dst.VersionNumber = src.Version
	dst.Status = schema.VersionStatus(src.Status)
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
}

func (repo *schemaVersionRepo) Draft(ctx context.Context, toDraft schema.Version) (*schema.Version, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaVersionRepo.Draft",
		trace.WithAttributes(
			attribute.String("version.schema_id", toDraft.SchemaID.String()),
			attribute.Int("version.version", toDraft.VersionNumber),
		),
	)
	defer span.End()

	sqlcSchemaVersion, err := repo.queries(ctx).DraftSchemaVersion(ctx, sqlc.DraftSchemaVersionParams{
		SchemaID: toDraft.SchemaID,
		Version:  toDraft.VersionNumber,
	})
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	span.SetAttributes(attribute.String("version.id", sqlcSchemaVersion.ID.String()))

	var createdSchemaVersion schema.Version
	mapSchemaVersionFromDB(&createdSchemaVersion, &sqlcSchemaVersion)
	return &createdSchemaVersion, nil
}

func (repo *schemaVersionRepo) Publish(ctx context.Context, toPublish schema.Version) error {
	ctx, span := repo.tracer.Start(ctx, "SchemaVersionsRepo.Publish",
		trace.WithAttributes(
			attribute.String("version.id", toPublish.ID.String()),
			attribute.String("version.schema_id", toPublish.SchemaID.String()),
		),
	)
	defer span.End()

	if err := repo.queries(ctx).PublishSchemaVersion(ctx, sqlc.PublishSchemaVersionParams{
		ID:       toPublish.ID,
		SchemaID: toPublish.SchemaID,
	}); err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}

	return nil
}

func (repo *schemaVersionRepo) Archive(ctx context.Context, toArchive schema.Version) error {
	ctx, span := repo.tracer.Start(ctx, "SchemaVersionsRepo.Archive",
		trace.WithAttributes(
			attribute.String("version.id", toArchive.ID.String()),
			attribute.String("version.schema_id", toArchive.SchemaID.String()),
		),
	)
	defer span.End()

	if err := repo.queries(ctx).ArchiveSchemaVersion(ctx, sqlc.ArchiveSchemaVersionParams{
		ID:       toArchive.ID,
		SchemaID: toArchive.SchemaID,
	}); err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}

	return nil
}

func (repo *schemaVersionRepo) GetLatest(ctx context.Context, schemaID uuid.UUID) (*schema.Version, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaVersionsRepo.GetLatest",
		trace.WithAttributes(
			attribute.String("version.schema_id", schemaID.String()),
		),
	)
	defer span.End()

	latest, err := repo.queries(ctx).GetLatestSchemaVersion(ctx, schemaID)
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	var found schema.Version
	mapSchemaVersionFromDB(&found, &latest)
	return &found, nil
}
