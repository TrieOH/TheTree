package badges

import (
	"context"
	"lib/telemetry"
	"time"
	"univents/models"

	"github.com/google/uuid"
)

// ListByUser returns a user's active badges grouped by origin (attendant,
// staff), each split into current (edition ends_at >= now) and past editions,
// most current first. Public: no identity or role required.
func (o *Operations) ListByUser(ctx context.Context, userID uuid.UUID) (*models.BadgeProfileGroups, error) {
	ctx, span := telemetry.StartSpan(ctx, "BadgesService.ListByUser")
	defer span.End()

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
