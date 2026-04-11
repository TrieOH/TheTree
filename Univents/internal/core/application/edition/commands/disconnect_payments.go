package commands

import (
	"context"
	"errors"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"

	"github.com/google/uuid"
)

func (uc *CommandService) DisconnectPayments(ctx context.Context, editionID uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "EditionService.DisconnectPayments")
	defer span.End()

	var sub *authz.UserSubject
	sub, err = authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	var edition *domain.Edition
	edition, err = uc.editions.GetByID(ctx, editionID)
	if err != nil {
		return err
	}

	if err = authz.Require(ctx, uc.az,
		authz.Subject("user", sub.ID),
		authz.Permission("disconnect_payments"),
		authz.Resource("event", edition.ID.String()),
	); err != nil {
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
