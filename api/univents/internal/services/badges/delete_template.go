package badges

import (
	"context"
	"lib/telemetry"
	"univents/models"

	idx "sdk/identityx"

	"github.com/google/uuid"
)

func (o *Operations) DeleteTemplate(ctx context.Context, templateID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "BadgeTemplatesCommands.DeleteTemplate")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return err
	}

	template, err := o.repo.GetByID(ctx, templateID)
	if err != nil {
		return err
	}

	edition, err := o.editions.GetByID(ctx, template.EditionID)
	if err != nil {
		return err
	}

	err = o.authz.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.repo.Delete(ctx, templateID)
}
