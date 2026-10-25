package service

import (
	"context"

	"Informd/models"
)

// TODO: kill this duplicated namespaced route — authz.Form already anchors via the form's namespace.
func (o *Operations) CreateNamespaced(ctx context.Context, payload models.CreateNamespacedStepFieldInput) (*models.Field, error) {
	_, err := o.Authz.Form(ctx, payload.FormID, models.FormMemberRoleAdmin)
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
	err = o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		created, err = o.Fields.Create(ctx, *field)
		if err != nil {
			return err
		}
		if payload.Type == models.FieldTypeSelect && payload.SelectConfig != nil {
			_, err = o.Fields.CreateSelectConfig(ctx, models.FieldSelectConfig{
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
