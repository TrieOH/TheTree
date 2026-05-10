package security

import (
	"IdentityX/contracts"
	"IdentityX/internal/platform/database/sqlc"
	"IdentityX/internal/shared/ports"
	"context"
	"lib/database"
	"lib/errx"
	"lib/xslices"

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
	dbe    *errx.DBHandler
}

var _ ports.KeysRepository = (*keyRepo)(nil)

func NewKeysRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer, dbe *errx.DBHandler) ports.KeysRepository {
	return &keyRepo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    dbe,
	}
}

func (repo *keyRepo) queries(ctx context.Context) *sqlc.Queries {
	if tx, ok := ctx.Value(database.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return repo.q.WithTx(tx)
	}
	return repo.q
}

func (repo *keyRepo) keySpan(ctx context.Context, op string) (context.Context, trace.Span) {
	return repo.tracer.Start(ctx, "KeyRepo."+op)
}

func mapKeyPairFromDB(src sqlc.KeyPair) contracts.Pair {
	return contracts.Pair{
		ID:              src.ID,
		KID:             src.Kid,
		ProjectID:       src.ProjectID,
		KeyType:         contracts.KeyType(src.KeyType),
		Algorithm:       contracts.Algorithm(src.Algorithm),
		PublicKey:       src.PublicKey,
		PrivateKey:      src.PrivateKey,
		Usage:           contracts.Usage(src.Usage),
		Status:          contracts.Status(src.Status),
		CreatedAt:       src.CreatedAt,
		ExpiresAt:       src.ExpiresAt,
		VerifyExpiresAt: src.VerifyExpiresAt,
	}
}

func mapPublicKeyFromDB(src sqlc.ListPublicKeysRow) contracts.PublicKey {
	return contracts.PublicKey{
		KID:       src.Kid,
		Algorithm: contracts.Algorithm(src.Algorithm),
		PublicKey: src.PublicKey,
		CreatedAt: src.CreatedAt,
		ExpiresAt: src.ExpiresAt,
	}
}

func (repo *keyRepo) CreateKeyPair(ctx context.Context, pair contracts.Pair) (*contracts.Pair, error) {
	ctx, span := repo.keySpan(ctx, "CreateKey")
	span.SetAttributes(attribute.String("key.kid", pair.KID))
	span.SetAttributes(attribute.String("key.type", string(pair.KeyType)))
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
		return nil, repo.dbe.DB(err, "key pair")
	}
	return new(mapKeyPairFromDB(row)), nil
}

func (repo *keyRepo) GetKeyByKID(ctx context.Context, kid string) (*contracts.Pair, error) {
	ctx, span := repo.keySpan(ctx, "GetKeyByKID")
	span.SetAttributes(attribute.String("key.kid", kid))
	defer span.End()
	row, err := repo.queries(ctx).GetKeyByKID(ctx, kid)
	if err != nil {
		return nil, repo.dbe.DB(err, "key pair")
	}
	return new(mapKeyPairFromDB(row)), nil
}

func (repo *keyRepo) GetActiveSigningKey(ctx context.Context, projectID *uuid.UUID) (*contracts.Pair, error) {
	ctx, span := repo.keySpan(ctx, "GetActiveSigningKey")
	if projectID != nil {
		span.SetAttributes(attribute.String("project.id", projectID.String()))
	}
	defer span.End()
	row, err := repo.queries(ctx).GetActiveSigningKey(ctx, projectID)
	if err != nil {
		return nil, repo.dbe.DB(err, "signing key")
	}
	return new(mapKeyPairFromDB(row)), nil
}

func (repo *keyRepo) GetActiveSigningKID(ctx context.Context, projectID *uuid.UUID) (string, error) {
	ctx, span := repo.keySpan(ctx, "GetActiveSigningKID")
	defer span.End()
	if projectID != nil {
		span.SetAttributes(attribute.String("project.id", projectID.String()))
	}
	kid, err := repo.queries(ctx).GetActiveSigningKID(ctx, projectID)
	if err != nil {
		return "", repo.dbe.DB(err, "signing kid")
	}
	return kid, nil
}

func (repo *keyRepo) ListPublicKeys(ctx context.Context, projectID *uuid.UUID) ([]contracts.PublicKey, error) {
	ctx, span := repo.keySpan(ctx, "ListPublicKeys")
	defer span.End()
	if projectID != nil {
		span.SetAttributes(attribute.String("project.id", projectID.String()))
	}
	rows, err := repo.queries(ctx).ListPublicKeys(ctx, projectID)
	if err != nil {
		return nil, repo.dbe.DB(err, "keys")
	}
	return xslices.MapSlice(rows, mapPublicKeyFromDB), nil
}

func (repo *keyRepo) RotateSigningKeys(ctx context.Context, projectID *uuid.UUID) error {
	ctx, span := repo.keySpan(ctx, "RotateSigningKeys")
	defer span.End()
	if projectID != nil {
		span.SetAttributes(attribute.String("project.id", projectID.String()))
	}
	if err := repo.queries(ctx).RotateSigningKeys(ctx, projectID); err != nil {
		return repo.dbe.DB(err, "keys")
	}
	return nil
}

func (repo *keyRepo) RevokeKeyByKID(ctx context.Context, kid string) error {
	ctx, span := repo.keySpan(ctx, "RevokeKeyByKID")
	span.SetAttributes(attribute.String("key.kid", kid))
	defer span.End()
	if err := repo.queries(ctx).RevokeKeyByKID(ctx, kid); err != nil {
		return repo.dbe.DB(err, "key")
	}
	return nil
}

func (repo *keyRepo) DeleteExpiredRevokedKeys(ctx context.Context) error {
	ctx, span := repo.keySpan(ctx, "DeleteExpiredRevokedKeys")
	defer span.End()
	if err := repo.queries(ctx).DeleteExpiredRevokedKeys(ctx); err != nil {
		return repo.dbe.DB(err, "keys")
	}
	return nil
}
