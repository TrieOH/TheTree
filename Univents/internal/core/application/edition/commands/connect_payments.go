package commands

import (
	"context"
	"errors"
	"fmt"
	"univents/internal/core/domain"
	"univents/internal/shared/authz"
	"univents/internal/shared/errx"

	"github.com/google/uuid"
)

func (uc *CommandService) ConnectPayments(ctx context.Context, triePaymentsCredentialID, editionID uuid.UUID, triePaymentsProvider string) (err error) {
	ctx, span := uc.tracer.Start(ctx, "EditionService.ConnectPayments")
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
		return fmt.Errorf("error getting edition: %w", err)
	}

	var allowed bool
	allowed, err = ga.Authz.Check().User(sub.ID).
		Object("payments").
		Action("connect").
		Scope(edition.GoauthScopeID).
		Allowed(ctx)
	if err != nil {
		return err
	}
	if !allowed {
		return errx.Forbidden("edition").SetMessage("insufficient permissions")
	}

	if edition.TriePaymentsCredentialID != nil {
		return errors.New("payment account already connected")
	}

	if err = uc.editions.ConnectPaymentsAccount(ctx, editionID, triePaymentsCredentialID, triePaymentsProvider); err != nil {
		return err
	}

	return nil
}
