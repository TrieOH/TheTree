package signatures

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) CancelRequest(ctx context.Context, requestID uuid.UUID, reason *string) error {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesService.CancelRequest")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	request, err := o.requests.GetRequestByID(ctx, requestID)
	if err != nil {
		return err
	}

	edition, err := o.editions.GetByID(ctx, request.EditionID)
	if err != nil {
		return err
	}

	err = o.authz.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	if request.Status != models.SignatureRequestStatusPending {
		return fun.ErrConflict("signature request is no longer pending")
	}

	return o.requests.CancelRequest(ctx, requestID, reason)
}
