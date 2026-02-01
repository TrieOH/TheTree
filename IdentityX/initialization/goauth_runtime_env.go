package initialization

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/apierr"
	"GoAuth/internal/crypto"
	"GoAuth/internal/domain/scopes"
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
)

func SetupRuntimeEnv(db *pgxpool.Pool) {
	queries := sqlc.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := queries.GetActiveSigningKeyForGoAuth(ctx)
	if err != nil {
		if apierr.IsNotFound(apierr.FromSQLC(err)) {
			pub, priv, err := crypto.GenerateEd25519()
			if err != nil {
				log.Fatalf("failed to generate GoAuth key: %v", err)
			}
			defer zero(priv)

			encryptedPriv, err := crypto.Encrypt(priv)
			if err != nil {
				log.Fatalf("failed to encrypt GoAuth key: %v", err)
			}

			kid := "goauth:" + ulid.Make().String()
			expiresAt := time.Now().Add(7 * 24 * time.Hour)

			_, err = queries.CreateKeyPair(ctx, sqlc.CreateKeyPairParams{
				Kid:        kid,
				ProjectID:  nil,
				KeyType:    "goauth",
				Algorithm:  "EdDSA",
				PublicKey:  pub,
				PrivateKey: encryptedPriv,
				Usage:      "sign",
				Status:     "active",
				ExpiresAt:  expiresAt,
			})

			if err != nil {
				if apierr.IsUniqueViolation(err) {
					log.Println("GoAuth signing key already created by another instance")
				} else {
					log.Fatalf("failed to create GoAuth signing key: %v", err)
				}
			} else {
				log.Println("Created GoAuth signing key")
			}
		} else {
			log.Fatalf("failed checking GoAuth signing key: %v", err)
		}
	}

	_, err = queries.CreateScope(ctx, sqlc.CreateScopeParams{
		Type:       string(scopes.ScopeTypeGlobal),
		ProjectID:  nil,
		Name:       nil,
		ExternalID: nil,
	})
	if err != nil {
		if apierr.IsUniqueViolation(err) {
			log.Println("GoAuth Global scope already created by another instance")
		} else {
			log.Fatalf("Failed to create GoAuth Global scope: %v", err)
		}
	} else {
		log.Println("Created GoAuth Global scope")
	}
}

func tryRotateGoAuthKeys(ctx context.Context, q *sqlc.Queries) error {
	key, err := q.GetActiveSigningKeyForGoAuth(ctx)
	if err != nil {
		if apierr.IsNotFound(err) {
			// defensive: no signing key → create
			return createGoAuthKey(ctx, q)
		}
		return err
	}

	if time.Until(key.ExpiresAt) > 24*time.Hour {
		return nil
	}

	if err := q.RotateSigningKeysForGoAuth(ctx); err != nil {
		return err
	}

	return createGoAuthKey(ctx, q)
}

func createGoAuthKey(ctx context.Context, q *sqlc.Queries) error {
	pub, priv, err := crypto.GenerateEd25519()
	defer zero(priv)
	if err != nil {
		return err
	}

	encryptedPriv, err := crypto.Encrypt(priv)
	if err != nil {
		return err
	}

	kid := "goauth:" + ulid.Make().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	_, err = q.CreateKeyPair(ctx, sqlc.CreateKeyPairParams{
		Kid:        kid,
		ProjectID:  nil,
		KeyType:    "goauth",
		Algorithm:  "EdDSA",
		PublicKey:  pub,
		PrivateKey: encryptedPriv,
		Usage:      "sign",
		Status:     "active",
		ExpiresAt:  expiresAt,
	})

	if apierr.IsUniqueViolation(err) {
		return nil
	}
	return err
}

func tryRotateProjectKeys(ctx context.Context, q *sqlc.Queries) error {
	projects, err := q.ListProjectsWithActiveSigningKeys(ctx)
	if err != nil {
		return err
	}

	for _, projectID := range projects {
		key, err := q.GetActiveSigningKeyForProject(ctx, projectID)
		if err != nil {
			if apierr.IsNotFound(err) {
				_ = createProjectKey(ctx, q, *projectID)
				continue
			}
			return err
		}

		if time.Until(key.ExpiresAt) > 24*time.Hour {
			continue
		}

		if err := q.RotateSigningKeysForProject(ctx, projectID); err != nil {
			return err
		}

		_ = createProjectKey(ctx, q, *projectID)
	}

	return nil
}

func createProjectKey(ctx context.Context, q *sqlc.Queries, projectID uuid.UUID) error {
	pub, priv, err := crypto.GenerateEd25519()
	defer zero(priv)
	if err != nil {
		return err
	}

	encryptedPriv, err := crypto.Encrypt(priv)
	if err != nil {
		return err
	}

	kid := "project:" + projectID.String() + ":" + ulid.Make().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	_, err = q.CreateKeyPair(ctx, sqlc.CreateKeyPairParams{
		Kid:        kid,
		ProjectID:  &projectID,
		KeyType:    "project",
		Algorithm:  "EdDSA",
		PublicKey:  pub,
		PrivateKey: encryptedPriv,
		Usage:      "sign",
		Status:     "active",
		ExpiresAt:  expiresAt,
	})

	// Rely on DB uniqueness for safety in concurrent rotations
	if apierr.IsUniqueViolation(err) {
		return nil
	}

	return err
}

func queriesWithTx(ctx context.Context, q *sqlc.Queries) *sqlc.Queries {
	if tx, ok := ctx.Value(transactions.TxKeyValue).(pgx.Tx); ok && tx != nil {
		return q.WithTx(tx)
	}
	return q
}

func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
