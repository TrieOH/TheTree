package queries

import (
	"Informd/models"
	"context"
	"lib/authz"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (s *QueryService) ListFormMembers(ctx context.Context, namespaceID, formID uuid.UUID) ([]models.FormMember, error) {
	ctx, span := s.tracer.Start(ctx, "NamespaceService.ListFormMembers")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	namespace, err := s.namespaces.GetByID(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	if sub.ID != namespace.OwnerID {
		_, err = s.namespaces.GetMember(ctx, sub.ID, namespace.ID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return nil, err
		}
		if err != nil {
			_, err = s.forms.GetMember(ctx, sub.ID, formID)
			if err != nil && !fun.Is(err, fun.CodeNotFound) {
				return nil, err
			}
			if err != nil {
				return nil, fun.ErrForbidden("insufficient permissions")
			}
		}
	}

	members, err := s.forms.ListDirectMembers(ctx, formID)
	if err != nil {
		return nil, err
	}

	namespaceMembers, err := s.namespaces.ListMembers(ctx, namespace.ID)
	if err != nil {
		return nil, err
	}

	for _, m := range namespaceMembers {
		members = append(members, models.FormMember{
			UserID:  m.UserID,
			FormID:  formID,
			Role:    models.FormMemberRole(m.Role),
			AddedAt: m.AddedAt,
			AddedBy: m.AddedBy,
		})
	}

	return members, nil
}
