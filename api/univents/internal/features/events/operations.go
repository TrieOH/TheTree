// Package events is the events feature: its repo, handlers and operations.
// Build wires it.
package events

import (
	"context"
	"lib/database"
	"lib/objectstorage"
	"lib/telemetry"
	"log/slog"
	idx "sdk/identityx"
	"univents/errs"
	"univents/internal/authz"
	"univents/models"
	"univents/ports"
	"uuid"
)

type Deps struct {
	Events ports.EventRepo
	Obj    *objectstorage.Client
	IDX    *idx.Client
	Badges ports.BadgeStaffOps
	Authz  *authz.Service
	// Tx opens the transactions multi-step writes run in; threaded
	// from boot, no package-level runner.
	Tx database.TxRunner
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

func (o *Operations) AddMember(ctx context.Context, eventID uuid.UUID, payload models.AddEventMemberInput) (*models.EventMember, error) {
	event, err := o.Events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	actor, err := o.IDX.Actors.GetByEmail(ctx, payload.Email)
	if err != nil {
		return nil, err
	}

	member, err := o.Events.AddEventMember(ctx, eventID, actor.ID, payload.Role)
	if err != nil {
		return nil, err
	}

	err = o.Badges.AwardStaffBadgesForUser(ctx, eventID, actor.ID)
	if err != nil {
		telemetry.Log().Error("failed to award staff badges",
			slog.String("event_id", eventID.String()),
			slog.String("user_id", actor.ID.String()),
			telemetry.Err(err))
	}

	return member, nil
}

func (o *Operations) Create(ctx context.Context, payload models.CreateEventInput) (*models.Event, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	event := &models.Event{
		OwnerID:      ident.Sub.ID,
		FullName:     payload.FullName,
		Acronym:      payload.Acronym,
		Slug:         payload.Slug,
		Description:  payload.Description,
		Status:       models.EventStatusDraft,
		ContactEmail: payload.ContactEmail,
	}

	var created *models.Event
	err = o.Tx.WithinTx(ctx, func(ctx context.Context) error {
		created, err = o.Events.Create(ctx, event)
		if err != nil {
			return err
		}

		_, err = o.Events.AddEventMember(ctx, created.ID, ident.Sub.ID, models.EventMemberRoleOwner)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

func (o *Operations) Discontinue(ctx context.Context, eventID uuid.UUID) error {
	event, err := o.Events.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event.Status != models.EventStatusActive {
		return errs.ErrEventNotActive
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Events.Discontinue(ctx, eventID)
}

func (o *Operations) GetBySlug(ctx context.Context, slug string) (*models.Event, error) {
	return o.Events.GetBySlug(ctx, slug)
}

func (o *Operations) ListJoined(ctx context.Context) ([]models.Event, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	events, err := o.Events.ListJoined(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (o *Operations) ListMembers(ctx context.Context, eventID uuid.UUID) ([]models.EventMember, error) {
	event, err := o.Events.GetByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleStaff)
	if err != nil {
		return nil, err
	}

	members, err := o.Events.ListEventMembers(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (o *Operations) ListOwned(ctx context.Context) ([]models.Event, error) {
	ident, err := idx.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}

	events, err := o.Events.ListOwned(ctx, ident.Sub.ID)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (o *Operations) ListPublic(ctx context.Context) ([]models.Event, error) {
	events, err := o.Events.ListPublic(ctx)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (o *Operations) Patch(ctx context.Context, payload models.PatchEventInput) (*models.Event, error) {
	existing, err := o.Events.GetByID(ctx, payload.EventID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, existing.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	event := &models.Event{
		FullName:     payload.FullName,
		Acronym:      payload.Acronym,
		Slug:         payload.Slug,
		Description:  payload.Description,
		LogoURL:      payload.LogoURL,
		BannerURL:    payload.BannerURL,
		ContactEmail: payload.ContactEmail,
	}

	return o.Events.Patch(ctx, payload.EventID, event)
}

func (o *Operations) Publish(ctx context.Context, eventID uuid.UUID) error {
	event, err := o.Events.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	if event.Status != models.EventStatusDraft {
		return errs.ErrEventNotDraft
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Events.Publish(ctx, eventID)
}

func (o *Operations) RemoveMember(ctx context.Context, eventID uuid.UUID, payload models.RemoveMemberInput) error {
	event, err := o.Events.GetByID(ctx, eventID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Event(ctx, event.ID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	actor, err := o.IDX.Actors.GetByEmail(ctx, payload.Email)
	if err != nil {
		return err
	}

	err = o.Events.RemoveEventMember(ctx, eventID, actor.ID)
	if err != nil {
		return err
	}

	err = o.Badges.RevokeStaffBadgesForUser(ctx, eventID, actor.ID)
	if err != nil {
		telemetry.Log().Error("failed to revoke staff badges",
			slog.String("event_id", eventID.String()),
			slog.String("user_id", actor.ID.String()),
			telemetry.Err(err))
	}

	return nil
}
