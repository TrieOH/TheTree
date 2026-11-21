// Package fields is the fields feature: its repo, handlers and operations.
// Build wires it.
package fields

import (
	"Informd/errs"
	"Informd/internal/authz"
	"Informd/models"
	"Informd/ports"
	"context"
	"lib/database"
	"lib/validator"
	"lib/xslices"
	"uuid"
)

type Deps struct {
	Forms      ports.FormsRepo
	Steps      ports.StepRepo
	Fields     ports.FieldsRepo
	Namespaces ports.NamespaceRepo
	Authz      *authz.Service
	// Tx opens the transactions multi-step writes run in; threaded from
	// boot, no package-level runner.
	Tx database.TxRunner
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

func (o *Operations) BulkEdit(ctx context.Context, formID uuid.UUID, payload []models.UpdateStepFieldInput) error {
	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	for _, p := range payload {
		if p.StepID != payload[0].StepID {
			return errs.ErrFieldsAcrossSteps
		}
	}

	fields := xslices.MapSlice(payload, models.UpdateStepFieldInputToField)
	return o.Fields.BulkEdit(ctx, fields)
}

func (o *Operations) Create(ctx context.Context, payload models.CreateStepFieldInput) (*models.Field, error) {
	_, err := o.Authz.Form(ctx, payload.FormID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	field := &models.Field{
		StepID:       payload.StepID,
		Key:          payload.Key,
		Title:        payload.Title,
		Description:  payload.Description,
		PositionHint: payload.PositionHint,
		Required:     payload.Required,
		Type:         payload.Type,
		Placeholder:  payload.Placeholder,
		DefaultValue: payload.DefaultValue,
		Config:       payload.Config,
	}
	err = validator.Validate(field)
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

func (o *Operations) Delete(ctx context.Context, formID, fieldID uuid.UUID) error {
	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Fields.Delete(ctx, fieldID)
}

func (o *Operations) EditSelectConfig(ctx context.Context, formID uuid.UUID, payload models.FieldSelectConfig) (*models.FieldSelectConfig, error) {
	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	return o.Fields.UpdateSelectConfig(ctx, payload)
}

func (o *Operations) GetSelectConfig(ctx context.Context, formID, fieldID uuid.UUID) (*models.FieldSelectConfig, error) {
	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Fields.GetSelectConfig(ctx, fieldID)
}

func (o *Operations) List(ctx context.Context, formID, stepID uuid.UUID) ([]models.Field, error) {
	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleMember)
	if err != nil {
		return nil, err
	}

	return o.Fields.ListByStepID(ctx, stepID)
}
