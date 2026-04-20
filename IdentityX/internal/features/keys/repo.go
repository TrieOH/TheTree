package keys

import (
	"IdentityX/internal/platform/database"
	"IdentityX/internal/platform/database/sqlc"
	"IdentityX/internal/shared/contracts"
	"IdentityX/internal/shared/ports"
	"context"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type keyRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
}

var _ ports.KeysRepository = (*keyRepo)(nil)

func NewRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.KeysRepository {
	return &keyRepo{
		q:      q,
		log:    log,
		tracer: tracer,
	}
}

func (repo *keyRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func mapKeyPairFromDB(dst *contracts.Pair, src *sqlc.KeyPair) {
	dst.ID = src.ID
	dst.KID = src.Kid
	dst.ProjectID = src.ProjectID
	dst.KeyType = contracts.Type(src.KeyType)
	dst.Algorithm = contracts.Algorithm(src.Algorithm)
	dst.PublicKey = src.PublicKey
	dst.PrivateKey = src.PrivateKey
	dst.Usage = contracts.Usage(src.Usage)
	dst.Status = contracts.Status(src.Status)
	dst.CreatedAt = src.CreatedAt
	dst.ExpiresAt = src.ExpiresAt
	dst.VerifyExpiresAt = src.VerifyExpiresAt
}

func mapGoAuthPublicKeyFromDB(dst *contracts.PublicKey, src *sqlc.ListActivePublicKeysForGoAuthRow) {
	dst.KID = src.Kid
	dst.Algorithm = contracts.Algorithm(src.Algorithm)
	dst.PublicKey = src.PublicKey
	dst.CreatedAt = src.CreatedAt
	dst.ExpiresAt = src.ExpiresAt
}

func mapProjectPublicKeyFromDB(dst *contracts.PublicKey, src *sqlc.ListActivePublicKeysForProjectRow) {
	dst.KID = src.Kid
	dst.Algorithm = contracts.Algorithm(src.Algorithm)
	dst.PublicKey = src.PublicKey
	dst.CreatedAt = src.CreatedAt
	dst.ExpiresAt = src.ExpiresAt
}

func (repo *keyRepo) CreateKeyPair(ctx context.Context, pair contracts.Pair) (*contracts.Pair, error) {
	ctx, span := repo.tracer.Start(ctx, "KeyRepo.CreateKey",
		trace.WithAttributes(
			attribute.String("key.kid", pair.KID),
			attribute.String("key.type", string(pair.KeyType)),
		),
	)
	defer span.End()

	row, err := repo.queries(ctx).CreateKeyPair(ctx, sqlc.CreateKeyPairParams{
		Kid:             pair.KID,
		ProjectID:       pair.ProjectID,
		KeyType:         string(pair.KeyType),
		Algorithm:       string(pair.Algorithm),
		PublicKey:       pair.PublicKey,
		PrivateKey:      pair.PrivateKey,
		Usage:           string(pair.Usage),
		Status:          string(pair.Status),
		ExpiresAt:       pair.ExpiresAt,
		VerifyExpiresAt: pair.VerifyExpiresAt,
	})
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	mapKeyPairFromDB(&pair, &row)
	return &pair, nil
}

func (repo *keyRepo) RotateGoAuthSigningKeys(ctx context.Context) error {
	ctx, span := repo.tracer.Start(ctx, "KeyRepo.RotateGoAuthSigningKeys")
	defer span.End()

	if err := repo.queries(ctx).RotateSigningKeysForGoAuth(ctx); err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}

func (repo *keyRepo) RotateProjectSigningKeys(ctx context.Context, projectID uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "KeyRepo.RotateProjectSigningKeys",
		trace.WithAttributes(attribute.String("project.id", projectID.String())),
	)
	defer span.End()

	if err := repo.queries(ctx).RotateSigningKeysForProject(ctx, &projectID); err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}

func (repo *keyRepo) GetActiveGoAuthSigningKey(ctx context.Context) (*contracts.Pair, error) {
	ctx, span := repo.tracer.Start(ctx, "KeyRepo.GetActiveGoAuthSigningKey")
	defer span.End()

	row, err := repo.queries(ctx).GetActiveSigningKeyForGoAuth(ctx)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	var pair contracts.Pair
	mapKeyPairFromDB(&pair, &row)
	return &pair, nil
}

func (repo *keyRepo) GetActiveProjectSigningKey(ctx context.Context, projectID uuid.UUID) (*contracts.Pair, error) {
	ctx, span := repo.tracer.Start(ctx, "KeyRepo.GetActiveProjectSigningKey",
		trace.WithAttributes(attribute.String("project.id", projectID.String())),
	)
	defer span.End()

	row, err := repo.queries(ctx).GetActiveSigningKeyForProject(ctx, &projectID)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	var pair contracts.Pair
	mapKeyPairFromDB(&pair, &row)
	return &pair, nil
}

func (repo *keyRepo) GetGoAuthKeyByKID(ctx context.Context, kid string) (*contracts.Pair, error) {
	ctx, span := repo.tracer.Start(ctx, "KeyRepo.GetGoAuthKeyByKID",
		trace.WithAttributes(attribute.String("key.kid", kid)),
	)
	defer span.End()

	row, err := repo.queries(ctx).GetGoAuthKeyByKID(ctx, kid)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	var pair contracts.Pair
	mapKeyPairFromDB(&pair, &row)
	return &pair, nil
}

func (repo *keyRepo) GetProjectKeyByKID(ctx context.Context, kid string) (*contracts.Pair, error) {
	ctx, span := repo.tracer.Start(ctx, "KeyRepo.GetProjectKeyByKID",
		trace.WithAttributes(attribute.String("key.kid", kid)),
	)
	defer span.End()

	row, err := repo.queries(ctx).GetProjectKeyByKID(ctx, kid)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	var pair contracts.Pair
	mapKeyPairFromDB(&pair, &row)
	return &pair, nil
}

func (repo *keyRepo) ListGoAuthPublicKeys(ctx context.Context) ([]contracts.PublicKey, error) {
	ctx, span := repo.tracer.Start(ctx, "KeyRepo.ListGoAuthPublicKeys")
	defer span.End()

	rows, err := repo.queries(ctx).ListActivePublicKeysForGoAuth(ctx)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	keys := make([]contracts.PublicKey, 0, len(rows))
	for _, row := range rows {
		var k contracts.PublicKey
		mapGoAuthPublicKeyFromDB(&k, &row)
		keys = append(keys, k)
	}

	return keys, nil
}

func (repo *keyRepo) ListProjectPublicKeys(ctx context.Context, projectID uuid.UUID) ([]contracts.PublicKey, error) {
	ctx, span := repo.tracer.Start(ctx, "KeyRepo.ListProjectPublicKeys",
		trace.WithAttributes(attribute.String("project.id", projectID.String())),
	)
	defer span.End()

	rows, err := repo.queries(ctx).ListActivePublicKeysForProject(ctx, &projectID)
	if err != nil {
		return nil, fail.From(err).RecordCtx(ctx)
	}

	keys := make([]contracts.PublicKey, 0, len(rows))
	for _, row := range rows {
		var k contracts.PublicKey
		mapProjectPublicKeyFromDB(&k, &row)
		keys = append(keys, k)
	}

	return keys, nil
}

func (repo *keyRepo) RevokeKeyByKID(ctx context.Context, kid string) error {
	ctx, span := repo.tracer.Start(ctx, "KeyRepo.RevokeKeyByKID",
		trace.WithAttributes(attribute.String("key.kid", kid)),
	)
	defer span.End()

	if err := repo.queries(ctx).RevokeKeyByKID(ctx, kid); err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}

func (repo *keyRepo) DeleteExpiredRevokedKeys(ctx context.Context) error {
	ctx, span := repo.tracer.Start(ctx, "KeyRepo.DeleteExpiredRevokedKeys")
	defer span.End()

	if err := repo.queries(ctx).DeleteExpiredRevokedKeys(ctx); err != nil {
		return fail.From(err).RecordCtx(ctx)
	}

	return nil
}

func (repo *keyRepo) GetActiveGoAuthSigningKID(ctx context.Context) (string, error) {
	ctx, span := repo.tracer.Start(ctx, "KeyRepo.GetActiveGoAuthSigningKID")
	defer span.End()

	kid, err := repo.queries(ctx).GetActiveGoAuthSigningKID(ctx)
	if err != nil {
		return "", fail.From(err).WithArgs("signing kid").RecordCtx(ctx)
	}

	return kid, nil
}

func (repo *keyRepo) GetActiveProjectSigningKID(ctx context.Context, projectID uuid.UUID) (string, error) {
	ctx, span := repo.tracer.Start(ctx, "KeyRepo.GetActiveProjectSigningKID",
		trace.WithAttributes(
			attribute.String("project.id", projectID.String()),
		),
	)
	defer span.End()

	kid, err := repo.queries(ctx).GetActiveProjectSigningKID(ctx, &projectID)
	if err != nil {
		return "", fail.From(err).WithArgs("signing kid").RecordCtx(ctx)
	}

	return kid, nil
}
