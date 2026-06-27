package commands

import (
	"context"
	idx "sdk/identityx"

	"Informd/models"

	"github.com/MintzyG/fun"
)

func (s *Command) CreateNamespaced(ctx context.Context, payload models.CreateNamespacedStepFieldInput) (*models.Field, error) {
	ctx, span := s.tracer.Start(ctx, "FieldService.CreateNamespaced")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	namespaceMember, err := s.namespaces.GetMember(ctx, ident.Sub.ID, payload.NamespaceID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return nil, err
	}
	if fun.Is(err, fun.CodeNotFound) {
		if namespaceMember.Role == models.NamespaceMemberRoleViewer {
			member, err := s.forms.GetMember(ctx, ident.Sub.ID, payload.FormID)
			if err != nil && !fun.Is(err, fun.CodeNotFound) {
				return nil, err
			}
			if err != nil {
				return nil, fun.ErrForbidden("insufficient permissions")
			}
			if member.Role == models.FormMemberRoleViewer {
				return nil, fun.ErrForbidden("insufficient permissions")
			}
		}
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
