package service

import (
	"context"
	"lib/telemetry"
	"univents/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) CancelRequest(ctx context.Context, requestID uuid.UUID, reason *string) error {
	ctx, span := telemetry.StartSpan(ctx, "SignaturesService.CancelRequest")
	defer span.End()

	request, err := o.Requests.GetRequestByID(ctx, requestID)
	if err != nil {
		return err
	}

	edition, err := o.Editions.GetByID(ctx, request.EditionID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	if request.Status != models.SignatureRequestStatusPending {
		return fun.ErrConflict("signature request is no longer pending")
	}

	return o.Requests.CancelRequest(ctx, requestID, reason)
}
