package service

import (
	"context"

	"Informd/models"

	"github.com/google/uuid"
)

// TODO: kill this duplicated namespaced route — authz.Form already anchors via the form's namespace.
func (o *Operations) GetFormResponseCount(ctx context.Context, _, formID uuid.UUID) (int, error) {
	_, err := o.Authz.Form(ctx, formID, models.FormMemberRoleMember)
	if err != nil {
		return 0, err
	}

	return o.Forms.ResponsesCount(ctx, formID)
}
