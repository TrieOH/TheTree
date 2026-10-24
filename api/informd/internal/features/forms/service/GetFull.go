package service

import (
	"context"

	"Informd/models"

	"github.com/google/uuid"
)

func (o *Operations) GetFull(ctx context.Context, formID uuid.UUID) (*models.FullForm, error) {
	form, err := o.Forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Form(ctx, form.ID, models.FormMemberRoleMember)
	if err != nil {
		return nil, err
	}

	steps, err := o.Steps.List(ctx, formID)
	if err != nil {
		return nil, err
	}
	fields, err := o.Fields.ListByFormID(ctx, formID)
	if err != nil {
		return nil, err
	}
	answers, err := o.Answers.GetByFormID(ctx, formID)
	if err != nil {
		return nil, err
	}
	responses, err := o.Responses.ListByForm(ctx, formID)
	if err != nil {
		return nil, err
	}
	responders, err := o.Responders.GetByFormID(ctx, formID)
	if err != nil {
		return nil, err
	}

	// index responders by ID
	responderByID := make(map[uuid.UUID]models.Responder, len(responders))
	for _, r := range responders {
		responderByID[r.ID] = r
	}

	// index responses by ID
	responseByID := make(map[uuid.UUID]models.Response, len(responses))
	for _, r := range responses {
		responseByID[r.ID] = r
	}

	// index answers by field ID
	answersByField := make(map[uuid.UUID][]models.FullAnswer)
	for _, a := range answers {
		email := "anonymous"
		if response, ok := responseByID[a.ResponseID]; ok {
			if response.ResponderID != nil {
				if r, ok := responderByID[*response.ResponderID]; ok {
					email = r.Email
				}
			} else if response.Email != nil {
				email = *response.Email
			}
		}
		if a.FieldID != nil {
			answersByField[*a.FieldID] = append(answersByField[*a.FieldID], models.FullAnswer{
				Answer:    a,
				Responder: email,
			})
		}
	}

	// index fields by step ID
	fieldsByStep := make(map[uuid.UUID][]models.FullField)
	for _, f := range fields {
		fieldsByStep[f.StepID] = append(fieldsByStep[f.StepID], models.FullField{
			Field:   f,
			Answers: answersByField[f.ID],
		})
	}

	// assemble steps
	fullSteps := make([]models.FullStep, len(steps))
	for i, s := range steps {
		fullSteps[i] = models.FullStep{
			Step:   s,
			Fields: fieldsByStep[s.ID],
		}
	}

	return &models.FullForm{
		Form:  *form,
		Steps: fullSteps,
	}, nil
}
