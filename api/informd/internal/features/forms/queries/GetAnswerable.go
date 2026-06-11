package queries

import (
	"Informd/models"
	"context"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (q *Queries) GetAnswerable(ctx context.Context, formID uuid.UUID) (*models.FormAnswerable, error) {
	ctx, span := q.tracer.Start(ctx, "FormService.GetAnswerable")
	defer span.End()

	form, err := q.forms.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if form.Status != models.FormStatusOpen {
		return nil, fun.ErrForbidden("form is not open for answers")
	}

	steps, err := q.steps.List(ctx, formID)
	if err != nil {
		return nil, err
	}
	fields, err := q.fields.ListByFormID(ctx, formID)
	if err != nil {
		return nil, err
	}

	// index fields by step ID
	fieldsByStep := make(map[uuid.UUID][]models.FieldAnswerable)
	for _, f := range fields {
		var selectConfig *models.FieldSelectConfig
		if f.Type == models.FieldTypeSelect {
			selectConfig, err = q.fields.GetSelectConfig(ctx, f.ID)
			if err != nil {
				return nil, err
			}
		}

		fieldsByStep[f.StepID] = append(fieldsByStep[f.StepID], models.FieldAnswerable{
			Field:             f,
			FieldSelectConfig: selectConfig,
		})
	}

	// assemble steps
	fullSteps := make([]models.StepAnswerable, len(steps))
	for i, s := range steps {
		fullSteps[i] = models.StepAnswerable{
			Step:   s,
			Fields: fieldsByStep[s.ID],
		}
	}

	return &models.FormAnswerable{
		Form:  *form,
		Steps: fullSteps,
	}, nil
}
