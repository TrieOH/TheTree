package security

import (
	"IdentityX/internal/platform/database"
	"IdentityX/internal/platform/database/sqlc"
	"IdentityX/internal/shared/contracts"
	"IdentityX/internal/shared/errx"
	"IdentityX/internal/shared/ports"
	"context"

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

func NewKeysRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.KeysRepository {
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
	dst.KeyType = contracts.KeyType(src.KeyType)
	dst.Algorithm = contracts.Algorithm(src.Algorithm)
	dst.PublicKey = src.PublicKey
	dst.PrivateKey = src.PrivateKey
	dst.Usage = contracts.Usage(src.Usage)
	dst.Status = contracts.Status(src.Status)
	dst.CreatedAt = src.CreatedAt
	dst.ExpiresAt = src.ExpiresAt
	dst.VerifyExpiresAt = src.VerifyExpiresAt
}

func mapPublicKeyFromDB(dst *contracts.PublicKey, src *sqlc.ListPublicKeysRow) {
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
		return nil, errx.DB(err, "key pair")
	}

	mapKeyPairFromDB(&pair, &row)
	return &pair, nil
}

func (repo *keyRepo) GetKeyByKID(ctx context.Context, kid string) (*contracts.Pair, error) {
	ctx, span := repo.tracer.Start(ctx, "KeyRepo.GetKeyByKID",
		trace.WithAttributes(attribute.String("key.kid", kid)),
	)
	defer span.End()

	row, err := repo.queries(ctx).GetKeyByKID(ctx, kid)
	if err != nil {
		return nil, errx.DB(err, "key pair")
	}

	var pair contracts.Pair
	mapKeyPairFromDB(&pair, &row)
	return &pair, nil
}

func (repo *keyRepo) GetActiveSigningKey(ctx context.Context, projectID *uuid.UUID) (*contracts.Pair, error) {
	var attrs []attribute.KeyValue
	if projectID != nil {
		attrs = append(attrs, attribute.String("project.id", projectID.String()))
	}

	ctx, span := repo.tracer.Start(ctx, "KeyRepo.GetActiveSigningKey", trace.WithAttributes(attrs...))
	defer span.End()

	row, err := repo.queries(ctx).GetActiveSigningKey(ctx, projectID)
	if err != nil {
		return nil, errx.DB(err, "signing key")
	}

	var pair contracts.Pair
	mapKeyPairFromDB(&pair, &row)
	return &pair, nil
}

func (repo *keyRepo) GetActiveSigningKID(ctx context.Context, projectID *uuid.UUID) (string, error) {
	var attrs []attribute.KeyValue
	if projectID != nil {
		attrs = append(attrs, attribute.String("project.id", projectID.String()))
	}

	ctx, span := repo.tracer.Start(ctx, "KeyRepo.GetActiveSigningKID", trace.WithAttributes(attrs...))
	defer span.End()

	kid, err := repo.queries(ctx).GetActiveSigningKID(ctx, projectID)
	if err != nil {
		return "", errx.DB(err, "signing kid")
	}

	return kid, nil
}

func (repo *keyRepo) ListPublicKeys(ctx context.Context, projectID *uuid.UUID) ([]contracts.PublicKey, error) {
	var attrs []attribute.KeyValue
	if projectID != nil {
		attrs = append(attrs, attribute.String("project.id", projectID.String()))
	}

	ctx, span := repo.tracer.Start(ctx, "KeyRepo.ListPublicKeys", trace.WithAttributes(attrs...))
	defer span.End()

	rows, err := repo.queries(ctx).ListPublicKeys(ctx, projectID)
	if err != nil {
		return nil, errx.DB(err, "keys")
	}

	keys := make([]contracts.PublicKey, 0, len(rows))
	for _, row := range rows {
		var k contracts.PublicKey
		mapPublicKeyFromDB(&k, &row)
		keys = append(keys, k)
	}

	return keys, nil
}

func (repo *keyRepo) RotateSigningKeys(ctx context.Context, projectID *uuid.UUID) error {
	var attrs []attribute.KeyValue
	if projectID != nil {
		attrs = append(attrs, attribute.String("project.id", projectID.String()))
	}

	ctx, span := repo.tracer.Start(ctx, "KeyRepo.RotateSigningKeys", trace.WithAttributes(attrs...))
	defer span.End()

	if err := repo.queries(ctx).RotateSigningKeys(ctx, projectID); err != nil {
		return errx.DB(err, "keys")
	}

	return nil
}

func (repo *keyRepo) RevokeKeyByKID(ctx context.Context, kid string) error {
	ctx, span := repo.tracer.Start(ctx, "KeyRepo.RevokeKeyByKID",
		trace.WithAttributes(attribute.String("key.kid", kid)),
	)
	defer span.End()

	if err := repo.queries(ctx).RevokeKeyByKID(ctx, kid); err != nil {
		return errx.DB(err, "key")
	}

	return nil
}

func (repo *keyRepo) DeleteExpiredRevokedKeys(ctx context.Context) error {
	ctx, span := repo.tracer.Start(ctx, "KeyRepo.DeleteExpiredRevokedKeys")
	defer span.End()

	if err := repo.queries(ctx).DeleteExpiredRevokedKeys(ctx); err != nil {
		return errx.DB(err, "keys")
	}

	return nil
}
