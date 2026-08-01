package fields

import (
	"context"
	"lib/telemetry"
	idx "sdk/identityx"

	"Informd/internal/authz"
	"Informd/models"

	"github.com/google/uuid"
)

func (o *Operations) Delete(ctx context.Context, formID, fieldID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "FieldService.Delete")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	err = authz.Service.CheckForm(ctx, ident.Sub.ID, formID, models.FormMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.fields.Delete(ctx, fieldID)
}
