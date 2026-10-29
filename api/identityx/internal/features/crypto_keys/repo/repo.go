package repo

import (
	"IdentityX/errs"
	"IdentityX/internal/sqlc"
	"IdentityX/models"
	"IdentityX/ports"
	"context"
	"errors"
	"lib/crypto"
	"lib/database"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repo struct {
	q   *sqlc.Queries
	dbe database.ErrorHandler
}

var _ ports.CryptoKeysRepo = (*Repo)(nil)

func New(q *sqlc.Queries) *Repo {
	return &Repo{
		q:   q,
		dbe: database.NewErrorHandler(errs.ErrCryptoKeyNotFound),
	}
}

// db is the queries handle, on the context's transaction when it carries one.
func (repo *Repo) db(ctx context.Context) *sqlc.Queries { return database.Queries(ctx, repo.q) }

func mapKeys(src sqlc.CryptoKey) models.CryptoKey {
	return models.CryptoKey{
		ID:                  src.ID,
		ProjectID:           src.ProjectID,
		Type:                models.CryptoKeyType(src.Type),
		Status:              models.CryptoKeyStatus(src.Status),
		PublicKey:           src.PublicKey,
		EncryptedPrivateKey: src.EncryptedPrivateKey,
		Algorithm:           src.Algorithm,
		Metadata:            src.Metadata,
		Active:              src.Active,
		CreatedAt:           src.CreatedAt,
		RotatedAt:           src.RotatedAt,
		ExpiresAt:           src.ExpiresAt,
	}
}

func mapToActiveSigningKey(src sqlc.GetActiveSigningKeysRow) models.ActiveSigningKey {
	return models.ActiveSigningKey{
		ID:        src.ID,
		PublicKey: src.PublicKey,
		Algorithm: src.Algorithm,
	}
}

func (repo *Repo) Create(ctx context.Context, projectID *uuid.UUID, pair *crypto.KeyPair, keyType string, expiresAt *time.Time) (*models.CryptoKey, error) {
	sqlcKeyPair, err := repo.db(ctx).CreateCryptoKey(ctx, sqlc.CreateCryptoKeyParams{
		ProjectID:           projectID,
		Type:                keyType,
		PublicKey:           pair.Public,
		EncryptedPrivateKey: pair.EncryptedPrivate,
		Algorithm:           pair.Algorithm,
		ExpiresAt:           expiresAt,
	})
	return database.One(sqlcKeyPair, err, repo.dbe, mapKeys)
}

func (repo *Repo) GetActive(ctx context.Context, keyType models.CryptoKeyType, projectID *uuid.UUID) (*models.CryptoKey, error) {
	sqlcKeyPair, err := repo.db(ctx).GetActiveCryptoKey(ctx, sqlc.GetActiveCryptoKeyParams{
		Type:      string(keyType),
		ProjectID: projectID,
	})
	return database.One(sqlcKeyPair, err, repo.dbe, mapKeys)
}

func (repo *Repo) GetActiveSigningKeys(ctx context.Context, projectID *uuid.UUID) ([]models.ActiveSigningKey, error) {
	sqlcKeys, err := repo.db(ctx).GetActiveSigningKeys(ctx, projectID)
	return database.Many(sqlcKeys, err, repo.dbe, mapToActiveSigningKey)
}

func (repo *Repo) GetByID(ctx context.Context, id uuid.UUID) (*models.CryptoKey, error) {
	sqlcKeyPair, err := repo.db(ctx).GetCryptoKeyByID(ctx, id)
	return database.One(sqlcKeyPair, err, repo.dbe, mapKeys)
}

// Retire conditionally moves the active key to retiring and stamps
// rotated_at. The WHERE status = 'active' guard makes the update
// conditional: when a concurrent rotation already retired the key, no row
// matches and Retire reports false, so the caller skips creating a
// replacement and overlapping workers cannot double-create.
func (repo *Repo) Retire(ctx context.Context, id uuid.UUID, rotatedAt time.Time) (bool, error) {
	_, err := repo.db(ctx).RetireCryptoKey(ctx, sqlc.RetireCryptoKeyParams{
		ID:        id,
		RotatedAt: &rotatedAt,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, repo.dbe(err)
	}
	return true, nil
}

// SweepRetiring retires the scope's keys whose grace period has elapsed:
// one RefreshTTL after rotation nothing valid is signed by them anymore.
func (repo *Repo) SweepRetiring(ctx context.Context, projectID *uuid.UUID, before time.Time) error {
	err := repo.db(ctx).SweepRetiringCryptoKeys(ctx, sqlc.SweepRetiringCryptoKeysParams{
		ProjectID: projectID,
		Before:    &before,
	})
	return repo.dbe(err)
}
