package queries

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (q *Queries) GetFullForm(ctx context.Context, namespaceID, formID uuid.UUID) (*models.FullForm, error) {
	ctx, span := q.tracer.Start(ctx, "NamespaceService.GetFullForm")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	namespace, err := q.namespaces.GetByID(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	form, err := q.forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if ident.Sub.ID != namespace.OwnerID {
		_, err = q.namespaces.GetMember(ctx, ident.Sub.ID, namespace.ID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return nil, err
		}
		if err != nil {
			_, err = q.forms.GetMember(ctx, ident.Sub.ID, formID)
			if err != nil && !fun.Is(err, fun.CodeNotFound) {
				return nil, err
			}
			if err != nil {
				return nil, fun.ErrForbidden("insufficient permissions")
			}
		}
	}

	steps, err := q.steps.List(ctx, formID)
	if err != nil {
		return nil, err
	}
	fields, err := q.fields.ListByFormID(ctx, formID)
	if err != nil {
		return nil, err
	}
	answers, err := q.answers.GetByFormID(ctx, formID)
	if err != nil {
		return nil, err
	}
	responses, err := q.responses.ListByForm(ctx, formID)
	if err != nil {
		return nil, err
	}
	responders, err := q.responders.GetByFormID(ctx, formID)
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
