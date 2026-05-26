package commands

import (
	"Informd/models"
	"context"
	"lib/authz"

	"github.com/MintzyG/fun"
)

func (s *Command) CreateNamespaced(ctx context.Context, payload models.CreateNamespacedFormStepInput) (*models.Step, error) {
	ctx, span := s.tracer.Start(ctx, "StepService.CreateNamespaced")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	step, err := models.NewStep(payload.FormID, payload.Title, payload.Description, payload.PositionHint)
	if err != nil {
		return nil, err
	}

	namespaceMember, err := s.namespaces.GetMember(ctx, sub.ID, payload.NamespaceID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return nil, err
	}
	if fun.Is(err, fun.CodeNotFound) {
		if namespaceMember.Role == models.NamespaceMemberRoleViewer {
			member, err := s.forms.GetMember(ctx, sub.ID, payload.FormID)
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

	created, err := s.steps.Create(ctx, *step)
	if err != nil {
		return nil, err
	}

	return created, nil
}
