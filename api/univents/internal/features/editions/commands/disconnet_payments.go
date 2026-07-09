package commands

import (
	"context"
	"errors"
	"univents/contracts"

	"github.com/google/uuid"
)

func (uc *Commands) DisconnectPayments(ctx context.Context, editionID uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "EditionService.DisconnectPayments")
	defer span.End()

	var edition *contracts.Edition
	edition, err = uc.editions.GetByID(ctx, editionID)
	if err != nil {
		return err
	}

	if edition.TriePaymentsCredentialID == nil {
		return errors.New("payment account already disconnected")
	}

	if err = uc.editions.DisconnectPaymentsAccount(ctx, editionID); err != nil {
		return err
	}

	return nil
}
