package commands

import (
	"context"
	"errors"
	"fmt"
	"univents/contracts"

	"github.com/google/uuid"
)

func (uc *Commands) ConnectPayments(ctx context.Context, triePaymentsCredentialID, editionID uuid.UUID, triePaymentsProvider, publicKey string) (err error) {
	ctx, span := uc.tracer.Start(ctx, "EditionService.ConnectPayments")
	defer span.End()

	var edition *contracts.Edition
	edition, err = uc.editions.GetByID(ctx, editionID)
	if err != nil {
		return fmt.Errorf("error getting edition: %w", err)
	}

	if edition.TriePaymentsCredentialID != nil {
		return errors.New("payment account already connected")
	}

	if err = uc.editions.ConnectPaymentsAccount(ctx, editionID, triePaymentsCredentialID, triePaymentsProvider, publicKey); err != nil {
		return err
	}

	return nil
}
