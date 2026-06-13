package commands

import (
	"context"

	"Informd/models"
	"lib/xslices"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (c *Commands) Submit(ctx context.Context, payload models.SubmitInput) error {
	ctx, span := c.tracer.Start(ctx, "ResponseService.Submit")
	defer span.End()

	form, err := c.forms.GetByID(ctx, payload.FormID)
	if err != nil {
		return err
	}

	if form.Status != models.FormStatusOpen {
		return fun.ErrBadRequest("form is not open for responses")
	}

	var responderID *uuid.UUID
	return c.tx.WithinTx(ctx, func(ctx context.Context) error {
		if payload.Email != nil {
			responder, err := c.responders.GetByEmail(ctx, *payload.Email)
			if err != nil && !fun.Is(err, fun.CodeNotFound) {
				return err
			}
			if fun.Is(err, fun.CodeNotFound) {
				responder, err = c.responders.Create(ctx, models.Responder{
					Email: *payload.Email,
				})
				if err != nil {
					return err
				}
			}
			responderID = &responder.ID
		}
		response, err := c.responses.Create(ctx, models.Response{
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

		err = c.answers.BatchUpsert(ctx, xslices.MapSlice(payload.Answers, models.SubmitAnswerInputToAnswer))
		if err != nil {
			return err
		}

		return c.responses.Finish(ctx, response.ID)
	})
}
