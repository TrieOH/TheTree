package badges

import (
	"context"
	"lib/telemetry"
	"strings"
	"time"
	"univents/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ListByUser returns a user's active badges grouped by origin (attendant,
// staff), each split into current (edition ends_at >= now) and past editions,
// most current first. Public: no identity or role required.
func (o *Operations) ListByUser(ctx context.Context, userID uuid.UUID) (*models.BadgeProfileGroups, error) {
	ctx, span := telemetry.StartSpan(ctx, "BadgesService.ListByUser")
	defer span.End()

	o.claimGiftsForUser(ctx, userID)

	views, err := o.emissions.ListViewsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	idx, err := o.loadTemplateIndex(ctx, views)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var attendantCurrent, attendantPast, staffCurrent, staffPast []profileBadgeEntry
	for _, v := range views {
		entry := profileBadgeEntry{badge: profileBadge(v, idx), endsAt: v.EndsAt}
		current := v.EndsAt.After(now) || v.EndsAt.Equal(now)
		switch v.Origin {
		case models.BadgeEmissionOriginStaff:
			if current {
				staffCurrent = append(staffCurrent, entry)
			} else {
				staffPast = append(staffPast, entry)
			}
		default:
			if current {
				attendantCurrent = append(attendantCurrent, entry)
			} else {
				attendantPast = append(attendantPast, entry)
			}
		}
	}

	return &models.BadgeProfileGroups{
		Attendant: flatten(attendantCurrent, attendantPast),
		Staff:     flatten(staffCurrent, staffPast),
	}, nil
}

// claimGiftsForUser ties any unclaimed email-only gifts for the user's
// account email to their id and emits the deferred badges, so a gifted
// badge appears the moment the recipient's profile loads. A profile
// requires an account, so the account email deterministically resolves the
// gifts — no edition-scoped UI visit needed (unlike the my-ticket claim).
// Best-effort: a claim or emit hiccup is logged and never fails the read
// (badges must render even if IdentityX is flaky).
func (o *Operations) claimGiftsForUser(ctx context.Context, userID uuid.UUID) {
	actor, err := o.actors.GetByID(ctx, userID)
	if err != nil {
		telemetry.Log().Warn("badge claim: actor lookup failed",
			zap.String("user_id", userID.String()), zap.Error(err))
		return
	}
	if actor.Email == nil {
		return
	}
	email := strings.TrimSpace(strings.ToLower(*actor.Email))
	if email == "" {
		return
	}

	claimed, err := o.registrations.ClaimAllByAttendeeEmail(ctx, email, userID)
	if err != nil {
		telemetry.Log().Warn("badge claim: gift claim failed",
			zap.String("user_id", userID.String()), zap.Error(err))
		return
	}
	for _, reg := range claimed {
		_, err = o.EmitForConfirmedRegistration(ctx, reg.ID)
		if err != nil {
			telemetry.Log().Warn("badge claim: emission failed",
				zap.String("registration_id", reg.ID.String()), zap.Error(err))
		}
	}
}
