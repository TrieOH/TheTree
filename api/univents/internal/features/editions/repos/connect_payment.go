package repos

import (
	"context"
	"lib/database"
	"lib/telemetry"
	"univents/internal/database/sqlc"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (repo *editionsRepo) ConnectPaymentsAccount(ctx context.Context, editionID, triePaymentsCredentialID uuid.UUID, triePaymentsProvider, publicKey string) error {
	ctx, span := repo.tracer.Start(ctx, "EditionsRepo.ConnectPaymentsAccount")
	defer span.End()

	telemetry.Log().Info("ConnectPaymentsAccount DATA", zap.String("edition_id", editionID.String()), zap.String("trie_payments_credential_id", triePaymentsCredentialID.String()), zap.String("trie_payments_provider", triePaymentsProvider))

	err := database.Queries(ctx, repo.q).ConnectEditionPaymentAccount(ctx, sqlc.ConnectEditionPaymentAccountParams{
		TriePaymentsCredentialID:      &triePaymentsCredentialID,
		TriePaymentsProvider:          &triePaymentsProvider,
		TriePaymentsProviderPublicKey: &publicKey,
		ID:                            editionID,
	})
	if err != nil {
		return repo.dbe(err)
	}

	return nil
}
