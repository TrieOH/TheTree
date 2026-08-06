package badges

import (
	"context"
	"lib/telemetry"
	"univents/models"

	idx "sdk/identityx"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

func (o *Operations) CreateTemplate(ctx context.Context, input models.CreateBadgeTemplateInput) (*models.BadgeTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "BadgeTemplatesCommands.CreateTemplate")
	defer span.End()

	err := validateTemplateScope(input)
	if err != nil {
		return nil, err
	}

	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	edition, err := o.editions.GetByID(ctx, input.EditionID)
	if err != nil {
		return nil, err
	}

	err = o.authz.CheckEvent(ctx, ident.Sub.ID, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	template := &models.BadgeTemplate{
		ID:           uuid.New(),
		EditionID:    input.EditionID,
		TicketTypeID: input.TicketTypeID,
		Origin:       input.Origin,
		Name:         input.Name,
		DesignData:   input.DesignData,
	}

	return o.repo.Create(ctx, template)
}

// validateTemplateScope keeps template scopes canonical: origin is either nil
// (default scope: edition default or per-ticket-type) or 'staff', and a staff
// template never targets a ticket type. The DB CHECK constraint
// (chk_badge_template_scope_valid) enforces the same.
func validateTemplateScope(input models.CreateBadgeTemplateInput) error {
	if input.Origin == nil {
		return nil
	}
	if *input.Origin != models.BadgeTemplateOriginStaff {
		return fun.ErrBadRequest("origin must be null or \"staff\"")
	}
	if input.TicketTypeID != nil {
		return fun.ErrBadRequest("a staff badge template cannot target a ticket type")
	}
	return nil
}
