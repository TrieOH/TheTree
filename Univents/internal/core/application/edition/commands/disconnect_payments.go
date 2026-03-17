package commands

import (
	"context"
	"errors"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
)

func (uc *CommandService) DisconnectPayments(ctx context.Context, editionID uuid.UUID) (err error) {
	ctx, span := uc.tracer.Start(ctx, "EditionService.DisconnectPayments")
	defer span.End()

	ga := uc.gaClient

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

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("payments").
		Action("disconnect").
		Scope(edition.GoauthScopeID).
		Allowed(ctx)
	if err != nil {
		return err
	}
	if !allowed {
		return errx.Forbidden("edition").SetMessage("insufficient permissions")
	}

	if edition.TriePaymentsCredentialID == nil {
		return errors.New("payment account already disconnected")
	}

	if err = uc.editions.DisconnectPaymentsAccount(ctx, editionID); err != nil {
		return err
	}

	return nil
}
