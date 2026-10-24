package service

import (
	"context"
	"lib/telemetry"
	"univents/models"
)

func (o *Operations) UpdateTemplate(ctx context.Context, input models.UpdateBadgeTemplateInput) (*models.BadgeTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "BadgesService.UpdateTemplate")
	defer span.End()

	template, err := o.Repo.GetByID(ctx, input.TemplateID)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, template.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
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
	return o.Repo.Update(ctx, template.ID, name, designData)
}
