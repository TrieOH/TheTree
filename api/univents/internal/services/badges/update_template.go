package badges

import (
	"context"
	"lib/telemetry"
	"univents/models"

	idx "sdk/identityx"
)

func (o *Operations) UpdateTemplate(ctx context.Context, input models.UpdateBadgeTemplateInput) (*models.BadgeTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "BadgesService.UpdateTemplate")
	defer span.End()

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	template, err := o.repo.GetByID(ctx, input.TemplateID)
	if err != nil {
		return nil, err
	}

	edition, err := o.editions.GetByID(ctx, template.EditionID)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	name := template.Name
	if input.Name != nil {
		name = *input.Name
	}
	designData := template.DesignData
	if input.DesignData != nil {
		designData = input.DesignData
	}

	// ticket_type_id is intentionally immutable: a template's target is set at
	// creation and changed by delete + create, keeping holder↔template mapping
	// consistent (designs are derived at read time anyway).
	return o.repo.Update(ctx, template.ID, name, designData)
}
