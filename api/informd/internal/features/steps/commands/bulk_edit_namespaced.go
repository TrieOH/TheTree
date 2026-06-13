package commands

import (
	"context"

	"Informd/models"
	"lib/authz"
	"lib/xslices"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *Command) BulkEditNamespaced(ctx context.Context, formID, namespaceID uuid.UUID, payload []models.UpdateNamespacedFormStepInput) error {
	ctx, span := s.tracer.Start(ctx, "StepService.BulkEditNamespaced")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return err
	}

	namespaceMember, err := s.namespaces.GetMember(ctx, sub.ID, namespaceID)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	if fun.Is(err, fun.CodeNotFound) {
		if namespaceMember.Role == models.NamespaceMemberRoleViewer {
			member, err := s.forms.GetMember(ctx, sub.ID, formID)
			if err != nil && !fun.Is(err, fun.CodeNotFound) {
				return err
			}
			if err != nil {
				return fun.ErrForbidden("insufficient permissions")
			}
			if member.Role == models.FormMemberRoleViewer {
				return fun.ErrForbidden("insufficient permissions")
			}
		}
	}

	for _, p := range payload {
		if p.FormID != formID {
			return fun.ErrBadRequest("all steps must belong to the same form")
		}
	}

	steps := xslices.MapSlice(payload, models.UpdateNamespacedFormStepInputToStep)
	return s.steps.BulkEdit(ctx, steps)
}
