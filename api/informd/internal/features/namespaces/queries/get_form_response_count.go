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
func (q *Queries) GetFormResponseCount(ctx context.Context, _, formID uuid.UUID) (int, error) {
	ctx, span := telemetry.StartSpan(ctx, "NamespaceService.GetFormResponseCount")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return 0, err
	}

	err = authz.Service.CheckForm(ctx, ident.Sub.ID, formID, models.FormMemberRoleMember)
	if err != nil {
		return 0, err
	}

	return q.forms.ResponsesCount(ctx, formID)
}
