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

type schemaVersionRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger // reserved for future use
	tracer trace.Tracer
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

func (r schemaVersionRepo) Draft(ctx context.Context, version schema.Version) (*schema.Version, error) {
	ctx, span := r.tracer.Start(ctx, "SchemaVersionRepo.Draft",
		trace.WithAttributes(
			attribute.String("version.schema_id", version.SchemaID.String()),
			attribute.Int("version.version", version.VersionNumber),
		),
	)
	defer span.End()

	sqlcSchemaVersion, err := r.q.DraftSchemaVersion(ctx, sqlc.DraftSchemaVersionParams{
		SchemaID: version.SchemaID,
		Version:  version.VersionNumber,
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

func (r schemaVersionRepo) Publish(ctx context.Context, version schema.Version) error {
	ctx, span := r.tracer.Start(ctx, "SchemaVersionsRepo.Publish",
		trace.WithAttributes(
			attribute.String("version.id", version.ID.String()),
			attribute.String("version.schema_id", version.SchemaID.String()),
		),
	)
	defer span.End()

	if err := r.q.PublishSchemaVersion(ctx, sqlc.PublishSchemaVersionParams{
		ID:       version.ID,
		SchemaID: version.SchemaID,
	}); err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}

	return nil
}

func (r schemaVersionRepo) Archive(ctx context.Context, version schema.Version) error {
	ctx, span := r.tracer.Start(ctx, "SchemaVersionsRepo.Archive",
		trace.WithAttributes(
			attribute.String("version.id", version.ID.String()),
			attribute.String("version.schema_id", version.SchemaID.String()),
		),
	)
	defer span.End()

	if err := r.q.ArchiveSchemaVersion(ctx, sqlc.ArchiveSchemaVersionParams{
		ID:       version.ID,
		SchemaID: version.SchemaID,
	}); err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return sqlcErr
	}

	return nil
}

func (r schemaVersionRepo) GetLatest(ctx context.Context, schemaID uuid.UUID) (*schema.Version, error) {
	ctx, span := r.tracer.Start(ctx, "SchemaVersionsRepo.Archive",
		trace.WithAttributes(
			attribute.String("version.schema_id", schemaID.String()),
		),
	)
	defer span.End()

	latest, err := r.q.GetLatestSchemaVersion(ctx, schemaID)
	if err != nil {
		sqlcErr := apierr.FromSQLC(err)
		apierr.RecordSQLCError(span, sqlcErr)
		return nil, sqlcErr
	}

	var found schema.Version
	mapSchemaVersionFromDB(&found, &latest)
	return &found, nil
}
