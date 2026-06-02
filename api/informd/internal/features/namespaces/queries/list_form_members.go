package queries

import (
	"Informd/models"
	"context"
	"lib/authz"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (q *Queries) ListFormMembers(ctx context.Context, namespaceID, formID uuid.UUID) ([]models.FormMember, error) {
	ctx, span := q.tracer.Start(ctx, "NamespaceService.ListFormMembers")
	defer span.End()

	sub, err := authz.RequireSubject(ctx)
	if err != nil {
		return nil, err
	}

	namespace, err := q.namespaces.GetByID(ctx, namespaceID)
	if err != nil {
		return nil, err
	}

	if sub.ID != namespace.OwnerID {
		_, err = q.namespaces.GetMember(ctx, sub.ID, namespace.ID)
		if err != nil && !fun.Is(err, fun.CodeNotFound) {
			return nil, err
		}
		if err != nil {
			_, err = q.forms.GetMember(ctx, sub.ID, formID)
			if err != nil && !fun.Is(err, fun.CodeNotFound) {
				return nil, err
			}
			if err != nil {
				return nil, fun.ErrForbidden("insufficient permissions")
			}
		}
	}

	members, err := q.forms.ListDirectMembers(ctx, formID)
	if err != nil {
		return nil, err
	}
	namespaceMembers, err := q.namespaces.ListMembers(ctx, namespace.ID)
	if err != nil {
		return nil, err
	}

	var formRoleRank = map[models.FormMemberRole]int{
		models.FormMemberRoleViewer: 0,
		models.FormMemberRoleEditor: 1,
		models.FormMemberRoleAdmin:  2,
	}

	// Index direct members by UserID for O(1) lookup during dedup.
	// Namespace membership wins unless the direct role is strictly higher.
	merged := make(map[uuid.UUID]models.FormMember, len(members)+len(namespaceMembers))
	for _, m := range members {
		merged[m.UserID] = m
	}
	for _, m := range namespaceMembers {
		ns := models.FormMember{
			UserID:  m.UserID,
			FormID:  formID,
			Role:    models.FormMemberRole(m.Role),
			AddedAt: m.AddedAt,
			AddedBy: m.AddedBy,
		}
		if existing, ok := merged[m.UserID]; !ok || formRoleRank[ns.Role] >= formRoleRank[existing.Role] {
			merged[m.UserID] = ns
		}
	}

	members = make([]models.FormMember, 0, len(merged))
	for _, m := range merged {
		members = append(members, m)
	}

	return members, nil
}
