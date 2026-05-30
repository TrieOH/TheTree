package commands

import (
	"Informd/models"
	"context"
	"lib/authz"

	"github.com/MintzyG/fun"
)

func (s *Command) Create(ctx context.Context, payload models.CreateStepFieldInput) (*models.Field, error) {
	ctx, span := s.tracer.Start(ctx, "FieldService.Create")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	member, err := s.forms.GetMember(ctx, sub.ID, payload.FormID)
	if err != nil {
		return nil, err
	}
	if member.Role == models.FormMemberRoleViewer {
		return nil, fun.ErrForbidden("insufficient permissions")
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
	if err = s.tx.WithinTx(ctx, func(ctx context.Context) error {
		created, err = s.fields.Create(ctx, *field)
		if err != nil {
			return err
		}
		if payload.Type == models.FieldTypeSelect && payload.SelectConfig != nil {
			_, err = s.fields.CreateSelectConfig(ctx, models.FieldSelectConfig{
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
	}); err != nil {
		return nil, err
	}
	return created, nil
}
