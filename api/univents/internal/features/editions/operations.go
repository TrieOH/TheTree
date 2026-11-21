// Package editions is the editions feature: its repo, handlers and
// operations. Build wires it.
package editions

import (
	"context"
	"lib/telemetry"
	"log/slog"
	"univents/errs"
	"univents/internal/authz"
	"univents/models"
	"univents/ports"
	"uuid"
)

type Deps struct {
	Events   ports.EventRepo
	Editions ports.EditionRepo
	Badges   ports.BadgeEditionOps
	Authz    *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

func (o *Operations) Create(ctx context.Context, payload models.CreateEditionInput) (*models.Edition, error) {
	event, err := o.Events.GetByID(ctx, payload.EventID)
	if err != nil {
		return nil, err
	}

	ident, err := o.Authz.Event(ctx, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	edition := &models.Edition{
		EventID:   payload.EventID,
		Name:      payload.Name,
		Slug:      payload.Slug,
		IsDraft:   true,
		StartsAt:  payload.StartsAt,
		EndsAt:    payload.EndsAt,
		CreatedBy: ident.Sub.ID,
	}

	return o.Editions.Create(ctx, edition)
}

func (o *Operations) GetActive(ctx context.Context, eventID uuid.UUID) (*models.Edition, error) {
	return o.Editions.GetActive(ctx, eventID)
}

func (o *Operations) GetByEventAndEditionSlug(ctx context.Context, eventSlug, editionSlug string) (*models.Edition, error) {
	event, err := o.Events.GetBySlug(ctx, eventSlug)
	if err != nil {
		return nil, err
	}

	return o.Editions.GetBySlug(ctx, event.ID, editionSlug)
}

func (o *Operations) GetPast(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	return o.Editions.GetPast(ctx, eventID)
}

func (o *Operations) GetUpcoming(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	return o.Editions.GetUpcoming(ctx, eventID)
}

func (o *Operations) ListDraft(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	event, err := o.Events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleStaff)
	if err != nil {
		return nil, err
	}

	return o.Editions.ListDraft(ctx, event.ID)
}

func (o *Operations) ListPublic(ctx context.Context, eventID uuid.UUID) ([]models.Edition, error) {
	event, err := o.Events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	return o.Editions.ListPublic(ctx, event.ID)
}

func (o *Operations) Patch(ctx context.Context, payload models.PatchEditionInput) (*models.Edition, error) {
	existing, err := o.Editions.GetByID(ctx, payload.EditionID)
	if err != nil {
		return nil, err
	}

	event, err := o.Events.GetByID(ctx, existing.EventID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	edition := &models.Edition{
		Name:                payload.Name,
		Slug:                payload.Slug,
		Tagline:             payload.Tagline,
		Description:         payload.Description,
		RegistrationOpensAt: payload.RegistrationOpensAt,
		StartsAt:            payload.StartsAt,
		EndsAt:              payload.EndsAt,
		LocationName:        payload.LocationName,
		LocationAddress:     payload.LocationAddress,
		LogoURL:             payload.LogoURL,
		BannerURL:           payload.BannerURL,
		ContactEmail:        payload.ContactEmail,
	}

	return o.Editions.Patch(ctx, existing.ID, edition)
}

func (o *Operations) Publish(ctx context.Context, editionID uuid.UUID) error {
	edition, err := o.Editions.GetByID(ctx, editionID)
	if err != nil {
		return err
	}

	if !edition.IsDraft {
		return errs.ErrEditionAlreadyPublished
	}

	event, err := o.Events.GetByID(ctx, edition.EventID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	err = o.Editions.Publish(ctx, editionID)
	if err != nil {
		return err
	}

	err = o.Badges.AwardStaffBadgesForEdition(ctx, editionID)
	if err != nil {
		telemetry.Log().ErrorContext(ctx, "failed to award staff badges for edition",
			slog.String("edition_id", editionID.String()),
			telemetry.Err(err))
	}

	return nil
}
