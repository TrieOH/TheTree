package commands

import (
	"Informd/models"
	"context"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *CommandService) isFormEditor(ctx context.Context, subID uuid.UUID, namespaceID *uuid.UUID, formID uuid.UUID) error {
	if err := s.isFormAdmin(ctx, subID, namespaceID, formID); err == nil {
		return nil
	}
	formMember, err := s.forms.GetMember(ctx, subID, formID)
	if fun.Is(err, fun.CodeNotFound) || err == nil && formMember.Role != models.FormMemberRoleEditor {
		return fun.ErrForbidden("insufficient permissions")
	} else if err != nil {
		return fun.ErrInternal(err.Error())
	}
	return nil
}
