package badges

import (
	"context"
	"sort"
	"time"
	"univents/models"

	"github.com/google/uuid"
)

// editionTemplates indexes one edition's templates for derivation:
// a specific ticket-type template wins over the edition default (NULL
// ticket_type_id); no match at all yields no template (placeholder badge).
type editionTemplates struct {
	defaultT *models.BadgeTemplate
	byTicket map[uuid.UUID]*models.BadgeTemplate
}

type templateIndex map[uuid.UUID]editionTemplates

func indexTemplates(templates []models.BadgeTemplate) templateIndex {
	idx := templateIndex{}
	for i := range templates {
		t := &templates[i]
		et := idx[t.EditionID]
		if et.byTicket == nil {
			et.byTicket = map[uuid.UUID]*models.BadgeTemplate{}
		}
		if t.TicketTypeID == nil {
			et.defaultT = t
		} else {
			et.byTicket[*t.TicketTypeID] = t
		}
		idx[t.EditionID] = et
	}
	return idx
}

func (idx templateIndex) match(editionID uuid.UUID, ticketTypeID *uuid.UUID) *models.BadgeTemplate {
	et, ok := idx[editionID]
	if !ok {
		return nil
	}
	if ticketTypeID != nil {
		if t, ok := et.byTicket[*ticketTypeID]; ok {
			return t
		}
	}
	return et.defaultT
}

func (o *Operations) loadTemplateIndex(ctx context.Context, views []models.BadgeEmissionView) (templateIndex, error) {
	editionIDs := make([]uuid.UUID, 0, len(views))
	seen := map[uuid.UUID]struct{}{}
	for _, v := range views {
		if _, ok := seen[v.EditionID]; ok {
			continue
		}
		seen[v.EditionID] = struct{}{}
		editionIDs = append(editionIDs, v.EditionID)
	}
	if len(editionIDs) == 0 {
		return templateIndex{}, nil
	}
	templates, err := o.repo.ListByEditionIDs(ctx, editionIDs)
	if err != nil {
		return nil, err
	}
	return indexTemplates(templates), nil
}

func profileBadge(v models.BadgeEmissionView, idx templateIndex) models.BadgeProfileBadge {
	t := idx.match(v.EditionID, v.TicketTypeID)
	badge := models.BadgeProfileBadge{
		EmissionID:  v.ID,
		EditionID:   v.EditionID,
		EditionName: v.EditionName,
		EventName:   v.EventName,
		Origin:      v.Origin,
		TicketName:  v.TicketName,
		ActionURL:   profileURL(v.UserID),
	}
	if t != nil {
		id := t.ID
		name := t.Name
		badge.TemplateID = &id
		badge.TemplateName = &name
		badge.DesignData = t.DesignData
	}
	return badge
}

type profileBadgeEntry struct {
	badge  models.BadgeProfileBadge
	endsAt time.Time
}

func flatten(current, past []profileBadgeEntry) models.BadgeOriginGroup {
	// Most current first: current editions end soonest first; past editions
	// are most recent first.
	sort.SliceStable(current, func(i, j int) bool { return current[i].endsAt.Before(current[j].endsAt) })
	sort.SliceStable(past, func(i, j int) bool { return past[i].endsAt.After(past[j].endsAt) })
	g := models.BadgeOriginGroup{
		Current: make([]models.BadgeProfileBadge, 0, len(current)),
		Past:    make([]models.BadgeProfileBadge, 0, len(past)),
	}
	for _, e := range current {
		g.Current = append(g.Current, e.badge)
	}
	for _, e := range past {
		g.Past = append(g.Past, e.badge)
	}
	return g
}
