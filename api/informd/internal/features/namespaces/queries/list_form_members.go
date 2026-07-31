package queries

import (
	"context"
	idx "sdk/identityx"

	"Informd/internal/authz"
	"Informd/models"
	"lib/telemetry"

	"github.com/google/uuid"
)

// TODO: kill this duplicated namespaced route — CheckForm already anchors via the form's namespace.
func (q *Queries) ListFormMembers(ctx context.Context, namespaceID, formID uuid.UUID) ([]models.FormMember, error) {
	ctx, span := telemetry.StartSpan(ctx, "NamespaceService.ListFormMembers")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	err = authz.Service.CheckForm(ctx, ident.Sub.ID, formID, models.FormMemberRoleMember)
	if err != nil {
		return nil, err
	}

	members, err := q.forms.ListDirectMembers(ctx, formID)
	if err != nil {
		return nil, err
	}
	namespaceMembers, err := q.namespaces.ListMembers(ctx, namespaceID)
	if err != nil {
		return nil, err
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
		if existing, ok := merged[m.UserID]; !ok || ns.Role.Rank() >= existing.Role.Rank() {
			merged[m.UserID] = ns
		}
	}

	members = make([]models.FormMember, 0, len(merged))
	for _, m := range merged {
		members = append(members, m)
	}

	return members, nil
}
