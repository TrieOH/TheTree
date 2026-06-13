package jobs

import (
	"IdentityX/internal/database/sqlc"
	"context"
	"lib/crypto"
	"time"

	"github.com/google/uuid"
	"github.com/riverqueue/river"
)

type CreateCryptoKeyArgs struct {
	ProjectID *uuid.UUID `json:"project_id,omitempty"`
	KeyType   string     `json:"key_type"`
}

func (CreateCryptoKeyArgs) Kind() string { return "crypto_key.create" }

type CreateCryptoKeyWorker struct {
	river.WorkerDefaults[CreateCryptoKeyArgs]
	queries *sqlc.Queries
}

func NewCreateCryptoKeyWorker(queries *sqlc.Queries) *CreateCryptoKeyWorker {
	return &CreateCryptoKeyWorker{queries: queries}
}

func (w *CreateCryptoKeyWorker) Work(ctx context.Context, job *river.Job[CreateCryptoKeyArgs]) error {
	key, err := crypto.GenerateKeyPair(job.Args.KeyType)
	if err != nil {
		return err
	}

	_, err = w.queries.CreateCryptoKey(ctx, sqlc.CreateCryptoKeyParams{
		ProjectID:           job.Args.ProjectID,
		Type:                job.Args.KeyType,
		PublicKey:           key.Public,
		EncryptedPrivateKey: key.EncryptedPrivate,
		Algorithm:           key.Algorithm,
		ExpiresAt:           new(time.Now().Add(7 * 24 * time.Hour)),
	})
	return err
}
