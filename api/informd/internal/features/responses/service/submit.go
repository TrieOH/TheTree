package service

import (
	"context"
	"errors"

	"Informd/errs"
	"Informd/models"
	"lib/database"
	"lib/xslices"

	"github.com/google/uuid"
)

func (o *Operations) Submit(ctx context.Context, payload models.SubmitInput) error {
	form, err := o.Forms.GetByID(ctx, payload.FormID)
	if err != nil {
		return err
	}

	if form.Status != models.FormStatusOpen {
		return errs.ErrFormNotOpenForResponses
	}

	var responderID *uuid.UUID
	return o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		if payload.Email != nil {
			responder, err := o.Responders.GetByEmail(ctx, *payload.Email)
			if err != nil && !errors.Is(err, database.ErrNotFound) {
				return err
			}
			if errors.Is(err, database.ErrNotFound) {
				responder, err = o.Responders.Create(ctx, models.Responder{
					Email: *payload.Email,
				})
				if err != nil {
					return err
				}
			}
			responderID = &responder.ID
		}
		response, err := o.Responses.Create(ctx, models.Response{
			FormID:      form.ID,
			InviteID:    nil,
			ResponderID: responderID,
			Email:       payload.Email,
		})
		if err != nil {
			return err
		}

		for i := range payload.Answers {
			payload.Answers[i].ResponseID = response.ID
		}

		err = o.Answers.BatchUpsert(ctx, xslices.MapSlice(payload.Answers, models.SubmitAnswerInputToAnswer))
		if err != nil {
			return err
		}

		return o.Responses.Finish(ctx, response.ID)
	})
}
