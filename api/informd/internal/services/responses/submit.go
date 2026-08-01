package responses

import (
	"context"

	"Informd/models"
	"lib/database"
	"lib/telemetry"
	"lib/xslices"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) Submit(ctx context.Context, payload models.SubmitInput) error {
	ctx, span := telemetry.StartSpan(ctx, "ResponseService.Submit")
	defer span.End()

	form, err := o.forms.GetByID(ctx, payload.FormID)
	if err != nil {
		return err
	}

	if form.Status != models.FormStatusOpen {
		return fun.ErrBadRequest("form is not open for responses")
	}

	var responderID *uuid.UUID
	return database.RunTx(ctx, func(ctx context.Context) error {
		if payload.Email != nil {
			responder, err := o.responders.GetByEmail(ctx, *payload.Email)
			if err != nil && !fun.Is(err, fun.CodeNotFound) {
				return err
			}
			if fun.Is(err, fun.CodeNotFound) {
				responder, err = o.responders.Create(ctx, models.Responder{
					Email: *payload.Email,
				})
				if err != nil {
					return err
				}
			}
			responderID = &responder.ID
		}
		response, err := o.responses.Create(ctx, models.Response{
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

		err = o.answers.BatchUpsert(ctx, xslices.MapSlice(payload.Answers, models.SubmitAnswerInputToAnswer))
		if err != nil {
			return err
		}

		return o.responses.Finish(ctx, response.ID)
	})
}
