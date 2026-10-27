package service

import (
	"context"
	"univents/errs"
	"univents/models"

	"github.com/google/uuid"
)

func (o *Operations) CancelRequest(ctx context.Context, requestID uuid.UUID, reason *string) error {
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
		return errs.ErrSignatureRequestNotPending
	}

	return o.Requests.CancelRequest(ctx, requestID, reason)
}
