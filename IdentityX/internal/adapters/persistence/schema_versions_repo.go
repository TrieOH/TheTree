package persistence

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/domain/version"
	"GoAuth/internal/errx"
	"GoAuth/internal/ports/outbounds"
	"context"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	if tx, ok := ctx.Value(transactions.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

var _ outbounds.SchemaVersionRepository = (*schemaVersionRepo)(nil)

func NewSchemaVersionRepo(q *sqlc.Queries, l *zap.Logger, tracer trace.Tracer) outbounds.SchemaVersionRepository {
	return &schemaVersionRepo{
		q:      q,
		log:    l,
		tracer: tracer,
	}
}

func mapSchemaVersionFromDB(dst *version.Version, src *sqlc.SchemaVersion) {
	dst.ID = src.ID
	dst.SchemaID = src.SchemaID
	dst.VersionNumber = src.Version
	dst.Status = version.Status(src.Status)
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
	dst.BasedOnVersionID = src.BasedOnVersionID
}

func (repo *schemaVersionRepo) Draft(ctx context.Context, toDraft version.Version) (*version.Version, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaVersionsRepo.Draft",
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
		return nil, fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(attribute.String("version.id", sqlcSchemaVersion.ID.String()))

	var createdSchemaVersion version.Version
	mapSchemaVersionFromDB(&createdSchemaVersion, &sqlcSchemaVersion)
	return &createdSchemaVersion, nil
}

func (repo *schemaVersionRepo) Publish(ctx context.Context, toPublish version.Version) error {
	ctx, span := repo.tracer.Start(ctx, "SchemaVersionsRepo.Publish",
		trace.WithAttributes(
			attribute.String("version.id", toPublish.ID.String()),
			attribute.String("version.schema_id", toPublish.SchemaID.String()),
		),
	)
	defer span.End()

	affectedRows, err := repo.queries(ctx).PublishSchemaVersion(ctx, sqlc.PublishSchemaVersionParams{
		ID:       toPublish.ID,
		SchemaID: toPublish.SchemaID,
	})
	if err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	if affectedRows == 0 {
		return fail.New(errx.SchemaVersionNotDraft).With(err)
	}

	return nil
}

func (repo *schemaVersionRepo) Archive(ctx context.Context, toArchive version.Version) error {
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
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}

func (repo *schemaVersionRepo) GetByID(ctx context.Context, versionID uuid.UUID) (*version.Version, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaVersionsRepo.GetByID",
		trace.WithAttributes(
			attribute.String("version.version_id", versionID.String()),
		),
	)
	defer span.End()

	sqlcVersion, err := repo.queries(ctx).GetVersionByID(ctx, versionID)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	var found version.Version
	mapSchemaVersionFromDB(&found, &sqlcVersion)
	return &found, nil
}

func (repo *schemaVersionRepo) GetCurrent(ctx context.Context, schemaID uuid.UUID) (*version.Version, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaVersionsRepo.GetCurrent",
		trace.WithAttributes(
			attribute.String("version.schema_id", schemaID.String()),
		),
	)
	defer span.End()

	sqlcVersion, err := repo.queries(ctx).GetCurrentSchemaVersion(ctx, schemaID)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	var currentSchemaVersion version.Version
	mapSchemaVersionFromDB(&currentSchemaVersion, &sqlcVersion)
	return &currentSchemaVersion, nil
}

func (repo *schemaVersionRepo) GetLatest(ctx context.Context, schemaID uuid.UUID) (*version.Version, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaVersionsRepo.GetLatest",
		trace.WithAttributes(
			attribute.String("version.schema_id", schemaID.String()),
		),
	)
	defer span.End()

	latest, err := repo.queries(ctx).GetLatestSchemaVersion(ctx, schemaID)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	var found version.Version
	mapSchemaVersionFromDB(&found, &latest)
	return &found, nil
}

func (repo *schemaVersionRepo) GetLatestForUpdate(ctx context.Context, schemaID uuid.UUID) (*version.Version, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaVersionsRepo.GetLatestForUpdate",
		trace.WithAttributes(
			attribute.String("version.schema_id", schemaID.String()),
		),
	)
	defer span.End()

	latest, err := repo.queries(ctx).GetLatestSchemaVersionForUpdate(ctx, schemaID)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	var found version.Version
	mapSchemaVersionFromDB(&found, &latest)
	return &found, nil
}

func (repo *schemaVersionRepo) List(ctx context.Context, schemaID uuid.UUID) ([]version.Version, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaVersionsRepo.List",
		trace.WithAttributes(
			attribute.String("schema.id", schemaID.String()),
		),
	)
	defer span.End()

	sqlcVersions, err := repo.queries(ctx).ListSchemaVersion(ctx, schemaID)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(attribute.Int("schema.count", len(sqlcVersions)))

	listed := make([]version.Version, 0, len(sqlcVersions))
	for _, v := range sqlcVersions {
		var found version.Version
		mapSchemaVersionFromDB(&found, &v)
		listed = append(listed, found)
	}

	return listed, nil
}

func (repo *schemaVersionRepo) CopyOnDraft(ctx context.Context, schemaVersionID uuid.UUID) (*version.Version, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaVersionsRepo.CopyOnDraft",
		trace.WithAttributes(
			attribute.String("version.id", schemaVersionID.String()),
		),
	)
	defer span.End()

	sqlcVersion, err := repo.queries(ctx).CopyVersionOnDraft(ctx, schemaVersionID)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	var copied version.Version
	mapSchemaVersionFromDB(&copied, &sqlcVersion)
	return &copied, nil
}

func (repo *schemaVersionRepo) GetByVersionNumber(ctx context.Context, schemaID uuid.UUID, versionNumber int) (*version.Version, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaVersionsRepo.GetByVersionNumber",
		trace.WithAttributes(
			attribute.String("version.id", schemaID.String()),
			attribute.Int("version.number", versionNumber),
		),
	)
	defer span.End()

	sqlcVersion, err := repo.queries(ctx).GetVersionByNumber(ctx, sqlc.GetVersionByNumberParams{
		SchemaID: schemaID,
		Version:  versionNumber,
	})
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	var found version.Version
	mapSchemaVersionFromDB(&found, &sqlcVersion)
	return &found, nil
}

func (repo *schemaVersionRepo) HasFields(ctx context.Context, versionID uuid.UUID) (bool, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaVersionsRepo.VersionHasFields",
		trace.WithAttributes(attribute.String("version_id", versionID.String())),
	)
	defer span.End()

	hasFields, err := repo.queries(ctx).VersionHasFields(ctx, versionID)
	if err != nil {
		return false, fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(attribute.Bool("has_fields", hasFields))
	return hasFields, nil
}

func (repo *schemaVersionRepo) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	ctx, span := repo.tracer.Start(ctx, "SchemaVersionsRepo.Exists",
		trace.WithAttributes(attribute.String("id", id.String())),
	)
	defer span.End()

	exists, err := repo.queries(ctx).SchemaVersionExists(ctx, id)
	if err != nil {
		return false, fail.From(err).RecordCtx(ctx)
	}

	span.SetAttributes(attribute.Bool("exists", exists))
	return exists, nil
}
