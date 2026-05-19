package commands

import (
	"Informd/models"
	"context"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *CommandService) isFormAdmin(ctx context.Context, subID uuid.UUID, namespaceID *uuid.UUID, formID uuid.UUID) error {
	if namespaceID != nil {
		member, err := s.namespaces.GetMember(ctx, subID, *namespaceID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return err
		}
		if err == nil && member.Role == models.NamespaceMemberRoleAdmin {
			return nil
		}
	}
	formMember, err := s.forms.GetMember(ctx, subID, formID)
	if fun.Is(err, fun.CodeNotFound) || err == nil && formMember.Role != models.FormMemberRoleAdmin {
		return fun.ErrForbidden("insufficient permissions")
	} else if err != nil {
		return fun.ErrInternal(err.Error())
	}
	return nil
}
