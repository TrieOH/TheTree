package fields

import (
	"context"
	"lib/database"
	"lib/telemetry"
	idx "sdk/identityx"

	"Informd/models"
)

func (o *Operations) Create(ctx context.Context, payload models.CreateStepFieldInput) (*models.Field, error) {
	ctx, span := telemetry.StartSpan(ctx, "FieldService.Create")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckForm(ctx, ident.Sub.ID, payload.FormID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	field, err := models.NewField(
		payload.StepID,
		payload.Key,
		payload.Title,
		payload.Description,
		payload.PositionHint,
		payload.Required,
		payload.Type,
		payload.Placeholder,
		payload.DefaultValue,
		payload.Config,
	)
	if err != nil {
		return nil, err
	}

	var created *models.Field
	err = database.RunTx(ctx, func(ctx context.Context) error {
		created, err = o.fields.Create(ctx, *field)
		if err != nil {
			return err
		}
		if payload.Type == models.FieldTypeSelect && payload.SelectConfig != nil {
			_, err = o.fields.CreateSelectConfig(ctx, models.FieldSelectConfig{
				FieldID:   created.ID,
				Behaviour: payload.SelectConfig.Behaviour,
				ValueType: payload.SelectConfig.ValueType,
				Options:   payload.SelectConfig.Options,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}
