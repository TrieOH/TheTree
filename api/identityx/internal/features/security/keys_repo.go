package security

import (
	"IdentityX/internal/database/sqlc"
	"IdentityX/internal/shared/ports"
	"IdentityX/models"
	"context"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type keyRepo struct {
	q      *sqlc.Queries
	log    *zap.Logger
	tracer trace.Tracer
	dbe    database.ErrorHandler
}

var _ ports.KeysRepository = (*keyRepo)(nil)

func NewKeysRepo(q *sqlc.Queries, log *zap.Logger, tracer trace.Tracer) ports.KeysRepository {
	return &keyRepo{
		q:      q,
		log:    log,
		tracer: tracer,
		dbe:    database.NewErrorHandler("key pair"),
	}
}

func mapKeyPairFromDB(src sqlc.KeyPair) models.Pair {
	return models.Pair{
		ID:              src.ID,
		KID:             src.Kid,
		ProjectID:       src.ProjectID,
		KeyType:         models.KeyType(src.KeyType),
		Algorithm:       models.Algorithm(src.Algorithm),
		PublicKey:       src.PublicKey,
		PrivateKey:      src.PrivateKey,
		Usage:           models.Usage(src.Usage),
		Status:          models.Status(src.Status),
		CreatedAt:       src.CreatedAt,
		ExpiresAt:       src.ExpiresAt,
		VerifyExpiresAt: src.VerifyExpiresAt,
	}
}

func mapPublicKeyFromDB(src sqlc.ListPublicKeysRow) models.PublicKey {
	return models.PublicKey{
		KID:       src.Kid,
		Algorithm: models.Algorithm(src.Algorithm),
		PublicKey: src.PublicKey,
		CreatedAt: src.CreatedAt,
		ExpiresAt: src.ExpiresAt,
	}
}

func (repo *keyRepo) CreateKeyPair(ctx context.Context, pair models.Pair) (*models.Pair, error) {
	ctx, span := repo.tracer.Start(ctx, "CreateKey")
	span.SetAttributes(attribute.String("key.kid", pair.KID))
	span.SetAttributes(attribute.String("key.type", string(pair.KeyType)))
	defer span.End()
	row, err := database.Queries(ctx, repo.q).CreateKeyPair(ctx, sqlc.CreateKeyPairParams{
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
		return nil, repo.dbe(err)
	}
	return new(mapKeyPairFromDB(row)), nil
}

func (repo *keyRepo) GetKeyByKID(ctx context.Context, kid string) (*models.Pair, error) {
	ctx, span := repo.tracer.Start(ctx, "GetKeyByKID")
	span.SetAttributes(attribute.String("key.kid", kid))
	defer span.End()
	row, err := database.Queries(ctx, repo.q).GetKeyByKID(ctx, kid)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return new(mapKeyPairFromDB(row)), nil
}

func (repo *keyRepo) GetActiveSigningKey(ctx context.Context, projectID *uuid.UUID) (*models.Pair, error) {
	ctx, span := repo.tracer.Start(ctx, "GetActiveSigningKey")
	if projectID != nil {
		span.SetAttributes(attribute.String("project.id", projectID.String()))
	}
	defer span.End()
	row, err := database.Queries(ctx, repo.q).GetActiveSigningKey(ctx, projectID)
	if err != nil {
		return nil, repo.dbe(err, "signing key")
	}
	return new(mapKeyPairFromDB(row)), nil
}

func (repo *keyRepo) GetActiveSigningKID(ctx context.Context, projectID *uuid.UUID) (string, error) {
	ctx, span := repo.tracer.Start(ctx, "GetActiveSigningKID")
	defer span.End()
	if projectID != nil {
		span.SetAttributes(attribute.String("project.id", projectID.String()))
	}
	kid, err := database.Queries(ctx, repo.q).GetActiveSigningKID(ctx, projectID)
	if err != nil {
		return "", repo.dbe(err, "signing kid")
	}
	return kid, nil
}

func (repo *keyRepo) ListPublicKeys(ctx context.Context, projectID *uuid.UUID) ([]models.PublicKey, error) {
	ctx, span := repo.tracer.Start(ctx, "ListPublicKeys")
	defer span.End()
	if projectID != nil {
		span.SetAttributes(attribute.String("project.id", projectID.String()))
	}
	rows, err := database.Queries(ctx, repo.q).ListPublicKeys(ctx, projectID)
	if err != nil {
		return nil, repo.dbe(err)
	}
	return xslices.MapSlice(rows, mapPublicKeyFromDB), nil
}

func (repo *keyRepo) RotateSigningKeys(ctx context.Context, projectID *uuid.UUID) error {
	ctx, span := repo.tracer.Start(ctx, "RotateSigningKeys")
	defer span.End()
	if projectID != nil {
		span.SetAttributes(attribute.String("project.id", projectID.String()))
	}
	if err := database.Queries(ctx, repo.q).RotateSigningKeys(ctx, projectID); err != nil {
		return repo.dbe(err)
	}
	return nil
}

func (repo *keyRepo) RevokeKeyByKID(ctx context.Context, kid string) error {
	ctx, span := repo.tracer.Start(ctx, "RevokeKeyByKID")
	span.SetAttributes(attribute.String("key.kid", kid))
	defer span.End()
	if err := database.Queries(ctx, repo.q).RevokeKeyByKID(ctx, kid); err != nil {
		return repo.dbe(err)
	}
	return nil
}

func (repo *keyRepo) DeleteExpiredRevokedKeys(ctx context.Context) error {
	ctx, span := repo.tracer.Start(ctx, "DeleteExpiredRevokedKeys")
	defer span.End()
	if err := database.Queries(ctx, repo.q).DeleteExpiredRevokedKeys(ctx); err != nil {
		return repo.dbe(err)
	}
	return nil
}
