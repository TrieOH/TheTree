package commands

import (
	"context"
	"errors"
	"fmt"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"

	"github.com/google/uuid"
)

func (uc *CommandService) ConnectPayments(ctx context.Context, triePaymentsCredentialID, editionID uuid.UUID, triePaymentsProvider, publicKey string) (err error) {
	ctx, span := uc.tracer.Start(ctx, "EditionService.ConnectPayments")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	var edition *domain.Edition
	edition, err = uc.editions.GetByID(ctx, editionID)
	if err != nil {
		return fmt.Errorf("error getting edition: %w", err)
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("connect_payments"),
		authz.Resource("edition", edition.ID.String()),
	); err != nil {
		return err
	}

	if edition.TriePaymentsCredentialID != nil {
		return errors.New("payment account already connected")
	}

	if err = uc.editions.ConnectPaymentsAccount(ctx, editionID, triePaymentsCredentialID, triePaymentsProvider, publicKey); err != nil {
		return err
	}

	return nil
}
