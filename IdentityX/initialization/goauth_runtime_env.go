package initialization

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/adapters/persistence/transactions"
	"GoAuth/internal/crypto"
	"GoAuth/internal/domain/scopes"
	"GoAuth/internal/errx"
	"context"
	"log"
	"time"

	"github.com/MintzyG/fail/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
)

func SetupRuntimeEnv(db *pgxpool.Pool) {
	queries := sqlc.New(db)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// First, rotate any expired keys to clear the way for new key creation
	if err := queries.RotateExpiredGoAuthKeys(ctx); err != nil {
		log.Printf("Warning: failed to rotate expired GoAuth keys: %v", err)
	}
	if err := queries.RotateExpiredProjectKeys(ctx); err != nil {
		log.Printf("Warning: failed to rotate expired project keys: %v", err)
	}

	// Also run the full key rotation logic to create new keys for projects without active keys
	if err := tryRotateGoAuthKeys(ctx, queries); err != nil {
		log.Printf("Warning: failed to rotate goauth keys: %v", err)
	}

	if err := tryRotateProjectKeys(ctx, queries); err != nil {
		log.Printf("Warning: failed to rotate project keys: %v", err)
	}

	_, err := queries.GetActiveSigningKeyForGoAuth(ctx)
	if err != nil {
		if fail.Is(fail.From(err), errx.SQLNotFound) {
			var pub string
			var priv []byte
			pub, priv, err = crypto.GenerateEd25519()
			if err != nil {
				log.Fatalf("failed to generate GoAuth key: %v", err)
			}
			defer zero(priv)

			var encryptedPriv []byte
			encryptedPriv, err = crypto.Encrypt(priv)
			if err != nil {
				log.Fatalf("failed to encrypt GoAuth key: %v", err)
			}

			kid := "goauth:" + ulid.Make().String()
			expiresAt := time.Now().Add(7 * 24 * time.Hour)

			_, err = queries.CreateKeyPair(ctx, sqlc.CreateKeyPairParams{
				Kid:             kid,
				ProjectID:       nil,
				KeyType:         "goauth",
				Algorithm:       "EdDSA",
				PublicKey:       pub,
				PrivateKey:      encryptedPriv,
				Usage:           "sign",
				Status:          "active",
				ExpiresAt:       expiresAt,
				VerifyExpiresAt: expiresAt.Add(7 * 24 * time.Hour),
			})

			fe := fail.From(err)

			if fe != nil {
				if errx.IsUniqueViolation(fe) {
					log.Println("GoAuth signing key already created by another instance")
				} else {
					log.Fatalf("failed to create GoAuth signing key: %v", fe.Error())
				}
			} else {
				log.Println("Created GoAuth signing key")
			}
		} else {
			log.Fatalf("failed checking GoAuth signing key: %v", err.Error())
		}
	}

	_, err = queries.CreateScope(ctx, sqlc.CreateScopeParams{
		Type:       string(scopes.ScopeTypeGlobal),
		ProjectID:  nil,
		Name:       nil,
		ExternalID: nil,
		ParentID:   nil,
	})

	fe := fail.From(err)

	if fe != nil {
		if fail.Is(fe, errx.SCOPEOneGlobal) || errx.IsUniqueViolation(fe) {
			log.Println("GoAuth Global scope already created by another instance")
		} else {
			log.Fatalf("Failed to create GoAuth Global scope: %v", fe.Error())
		}
	} else {
		log.Println("Created GoAuth Global scope")
	}
}

func tryRotateGoAuthKeys(ctx context.Context, q *sqlc.Queries) error {
	key, err := q.GetActiveSigningKeyForGoAuth(ctx)
	if err != nil {
		if fail.Is(fail.From(err), errx.SQLNotFound) {
			// defensive: no signing key → create
			return createGoAuthKey(ctx, q)
		}
		return fail.From(err)
	}

	if time.Until(key.ExpiresAt) > 24*time.Hour {
		return nil
	}

	if err = q.RotateSigningKeysForGoAuth(ctx); err != nil {
		return fail.From(err)
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
		Kid:             kid,
		ProjectID:       nil,
		KeyType:         "goauth",
		Algorithm:       "EdDSA",
		PublicKey:       pub,
		PrivateKey:      encryptedPriv,
		Usage:           "sign",
		Status:          "active",
		ExpiresAt:       expiresAt,
		VerifyExpiresAt: expiresAt.Add(7 * 24 * time.Hour),
	})

	fe := fail.From(err)
	if errx.IsUniqueViolation(fe) {
		return nil
	}
	return fe
}

func tryRotateProjectKeys(ctx context.Context, q *sqlc.Queries) error {
	projects, err := q.ListProjectsWithSigningKeys(ctx)
	fe := fail.From(err)
	if fe != nil {
		return fe
	}

	for _, projectID := range projects {
		var key sqlc.KeyPair
		key, err = q.GetActiveSigningKeyForProject(ctx, projectID)
		if err != nil {
			if fail.Is(fail.From(err), errx.SQLNotFound) {
				_ = createProjectKey(ctx, q, *projectID)
				continue
			}
			return fail.From(err)
		}

		if time.Until(key.ExpiresAt) > 24*time.Hour {
			continue
		}

		if err = q.RotateSigningKeysForProject(ctx, projectID); err != nil {
			return fail.From(err)
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
		Kid:             kid,
		ProjectID:       &projectID,
		KeyType:         "project",
		Algorithm:       "EdDSA",
		PublicKey:       pub,
		PrivateKey:      encryptedPriv,
		Usage:           "sign",
		Status:          "active",
		ExpiresAt:       expiresAt,
		VerifyExpiresAt: expiresAt.Add(7 * 24 * time.Hour),
	})

	if err == nil {
		return nil
	}

	// Rely on DB uniqueness for safety in concurrent rotations
	if errx.IsUniqueViolation(fail.From(err)) {
		return nil
	}

	return fail.From(err)
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
