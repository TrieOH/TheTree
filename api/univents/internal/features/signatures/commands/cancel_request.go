package commands

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/internal/authz"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (c *Commands) CancelRequest(ctx context.Context, requestID uuid.UUID, reason *string) error {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesService.CancelRequest")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	request, err := c.requests.GetRequestByID(ctx, requestID)
	if err != nil {
		return err
	}

	edition, err := c.editions.GetByID(ctx, request.EditionID)
	if err != nil {
		return err
	}

	err = authz.Service.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	if request.Status != models.SignatureRequestStatusPending {
		return fun.ErrConflict("signature request is no longer pending")
	}

	return c.requests.CancelRequest(ctx, requestID, reason)
}
