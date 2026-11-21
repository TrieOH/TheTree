// Package badges is the badges feature: its repo, handlers and operations.
// Build wires it.
package badges

import (
	"context"
	"encoding/json"
	libjobs "lib/jobs"
	"lib/telemetry"
	"log/slog"
	"sort"
	"strings"
	"time"
	"univents/errs"
	"univents/internal/authz"
	"univents/internal/features/badges/jobs"
	"univents/models"
	"univents/ports"
	"uuid"
)

// Jobs is the background-work surface: the badge email is enqueued inside
// the caller's tx with a delay. Satisfied by the lib/jobs carriers.
type Jobs interface {
	EnqueueAfter(ctx context.Context, job libjobs.Job, delay time.Duration) error
}

type Deps struct {
	Repo          ports.BadgeTemplateRepo
	Emissions     ports.BadgeEmissionRepo
	Registrations ports.RegistrationRepo
	Editions      ports.EditionRepo
	Events        ports.EventRepo
	Actors        ports.ActorResolver
	Jobs          Jobs
	Authz         *authz.Service
	AppURL        string // the Univents front, for the profile links on badges
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

// AwardStaffBadgesForEdition emits staff badges for every member of the
// edition's event. Called when an edition is published so staff added before
// publication still get their badge.
func (o *Operations) AwardStaffBadgesForEdition(ctx context.Context, editionID uuid.UUID) error {
	edition, err := o.Editions.GetByID(ctx, editionID)
	if err != nil {
		return err
	}

	members, err := o.Events.ListEventMembers(ctx, edition.EventID)
	if err != nil {
		return err
	}

	for _, member := range members {
		_, err = o.Emissions.Upsert(ctx, &models.BadgeEmission{
			EditionID: edition.ID,
			UserID:    member.UserID,
			Origin:    models.BadgeEmissionOriginStaff,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// AwardStaffBadgesForUser emits a staff badge for the user on every published
// edition of the event that is current or future (ends_at >= now). Past
// editions are never awarded retroactively; draft editions are awarded when
// published (see AwardStaffBadgesForEdition).
func (o *Operations) AwardStaffBadgesForUser(ctx context.Context, eventID, userID uuid.UUID) error {
	editions, err := o.Editions.ListPublic(ctx, eventID)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, edition := range editions {
		if edition.EndsAt.Before(now) {
			continue
		}
		_, err = o.Emissions.Upsert(ctx, &models.BadgeEmission{
			EditionID: edition.ID,
			UserID:    userID,
			Origin:    models.BadgeEmissionOriginStaff,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *Operations) CreateTemplate(ctx context.Context, input models.CreateBadgeTemplateInput) (*models.BadgeTemplate, error) {
	err := validateTemplateScope(input)
	if err != nil {
		return nil, err
	}

	edition, err := o.Editions.GetByID(ctx, input.EditionID)
	if err != nil {
		return nil, err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
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

	return o.Repo.Create(ctx, template)
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
		return errs.ErrInvalidBadgeOrigin
	}
	if input.TicketTypeID != nil {
		return errs.ErrStaffTemplateWithTicketType
	}
	return nil
}

func (o *Operations) DeleteTemplate(ctx context.Context, templateID uuid.UUID) error {
	template, err := o.Repo.GetByID(ctx, templateID)
	if err != nil {
		return err
	}

	edition, err := o.Editions.GetByID(ctx, template.EditionID)
	if err != nil {
		return err
	}

	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return err
	}

	return o.Repo.Delete(ctx, templateID)
}

// badgeEmailDelay holds the badge email back long enough for the caller's tx
// to commit: it is enqueued inside the tx, before commit. On rollback the job
// finds no confirmed registration and no-ops.
const badgeEmailDelay = 30 * time.Second

// EmitForConfirmedRegistration emits the participant badge for a confirmed
// registration and enqueues its email (the badges.send_email job sends it
// exactly once). Called by the checkout feature when a registration is
// confirmed (paid: webhook approved; free: confirmed at inscription).
// Non-confirmed registrations are a no-op.
func (o *Operations) EmitForConfirmedRegistration(ctx context.Context, registrationID uuid.UUID) (*models.BadgeEmission, error) {
	reg, err := o.Registrations.GetByID(ctx, registrationID)
	if err != nil {
		return nil, err
	}
	if reg.Status != models.RegistrationStatusConfirmed {
		//nolint:nilnil // no-op contract: pending/failed registrations simply never emit
		return nil, nil
	}
	if reg.AttendeeUserID == nil {
		// Accountless gifted ticket (email-only recipient): there is no
		// profile to attach the badge to or email it to (the badge email's
		// QR points at the profile). Emission is deferred until the
		// recipient claims an account and the registration gets its actor
		// id — the claim flow re-emits from here.
		telemetry.Log().InfoContext(ctx, "badge emission deferred: attendee has no account yet",
			slog.String("registration_id", registrationID.String()))
		//nolint:nilnil // deferred emission is a normal state, not an error
		return nil, nil
	}

	emission, err := o.Emissions.Upsert(ctx, &models.BadgeEmission{
		EditionID:      reg.EditionID,
		UserID:         *reg.AttendeeUserID,
		Origin:         models.BadgeEmissionOriginParticipant,
		RegistrationID: &reg.ID,
	})
	if err != nil {
		return nil, err
	}

	if emission.EmailSentAt == nil {
		err = o.Jobs.EnqueueAfter(ctx, jobs.SendBadgeEmailArgs{EmissionID: emission.ID, RegistrationID: reg.ID}, badgeEmailDelay)
		if err != nil {
			telemetry.Log().ErrorContext(ctx, "failed to enqueue badge email",
				slog.String("registration_id", reg.ID.String()),
				telemetry.Err(err))
		}
	}
	return emission, nil
}

func (o *Operations) GetTemplate(ctx context.Context, templateID uuid.UUID) (*models.BadgeTemplate, error) {
	return o.Repo.GetByID(ctx, templateID)
}

// ListByEdition returns all emissions of an edition for the owner/admin to
// browse (all statuses; holder names are derived front-side).
func (o *Operations) ListByEdition(ctx context.Context, editionID uuid.UUID) ([]models.BadgeEditionEmission, error) {
	edition, err := o.Editions.GetByID(ctx, editionID)
	if err != nil {
		return nil, err
	}
	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	views, err := o.Emissions.ListViewsByEdition(ctx, editionID)
	if err != nil {
		return nil, err
	}
	idx, err := o.loadTemplateIndex(ctx, views)
	if err != nil {
		return nil, err
	}

	items := make([]models.BadgeEditionEmission, 0, len(views))
	for _, v := range views {
		item := models.BadgeEditionEmission{
			ID:           v.ID,
			UserID:       v.UserID,
			Origin:       v.Origin,
			Status:       v.Status,
			StatusReason: v.StatusReason,
			TicketName:   v.TicketName,
			EmittedAt:    v.EmittedAt,
			UpdatedAt:    v.UpdatedAt,
		}
		if t := idx.match(v.EditionID, v.TicketTypeID, v.Origin); t != nil {
			id := t.ID
			name := t.Name
			item.TemplateID = &id
			item.TemplateName = &name
		}
		items = append(items, item)
	}
	return items, nil
}

// ListByUser returns a user's active badges grouped by origin (attendant,
// staff), each split into current (edition ends_at >= now) and past editions,
// most current first. Public: no identity or role required.
func (o *Operations) ListByUser(ctx context.Context, userID uuid.UUID) (*models.BadgeProfileGroups, error) {
	o.claimGiftsForUser(ctx, userID)

	views, err := o.Emissions.ListViewsByUser(ctx, userID)
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
		entry := profileBadgeEntry{badge: o.profileBadge(v, idx), endsAt: v.EndsAt}
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
	actor, err := o.Actors.GetByID(ctx, userID)
	if err != nil {
		telemetry.Log().WarnContext(ctx, "badge claim: actor lookup failed",
			slog.String("user_id", userID.String()), telemetry.Err(err))
		return
	}
	if actor.Email == nil {
		return
	}
	email := strings.TrimSpace(strings.ToLower(*actor.Email))
	if email == "" {
		return
	}

	claimed, err := o.Registrations.ClaimAllByAttendeeEmail(ctx, email, userID)
	if err != nil {
		telemetry.Log().WarnContext(ctx, "badge claim: gift claim failed",
			slog.String("user_id", userID.String()), telemetry.Err(err))
		return
	}
	for _, reg := range claimed {
		_, err = o.EmitForConfirmedRegistration(ctx, reg.ID)
		if err != nil {
			telemetry.Log().WarnContext(ctx, "badge claim: emission failed",
				slog.String("registration_id", reg.ID.String()), telemetry.Err(err))
		}
	}
}

func (o *Operations) ListTemplates(ctx context.Context, editionID uuid.UUID) ([]models.BadgeTemplate, error) {
	return o.Repo.ListByEdition(ctx, editionID)
}

// PrintByEdition returns the print payload for the edition's active emissions,
// optionally filtered to specific emission ids. Owner/admin only.
func (o *Operations) PrintByEdition(ctx context.Context, editionID uuid.UUID, emissionIDs []uuid.UUID) ([]models.BadgePrintItem, error) {
	edition, err := o.Editions.GetByID(ctx, editionID)
	if err != nil {
		return nil, err
	}
	_, err = o.Authz.Event(ctx, edition.EventID, models.EventMemberRoleAdmin)
	if err != nil {
		return nil, err
	}

	views, err := o.Emissions.ListViewsByEdition(ctx, editionID)
	if err != nil {
		return nil, err
	}
	idx, err := o.loadTemplateIndex(ctx, views)
	if err != nil {
		return nil, err
	}

	include := map[uuid.UUID]struct{}{}
	for _, id := range emissionIDs {
		include[id] = struct{}{}
	}

	items := make([]models.BadgePrintItem, 0, len(views))
	for _, v := range views {
		if v.Status != models.BadgeEmissionStatusActive {
			continue
		}
		if len(include) > 0 {
			if _, ok := include[v.ID]; !ok {
				continue
			}
		}
		item := models.BadgePrintItem{
			EmissionID:  v.ID,
			UserID:      v.UserID,
			Origin:      v.Origin,
			EventName:   v.EventName,
			EditionName: v.EditionName,
			TicketName:  v.TicketName,
			ActionURL:   o.profileURL(v.UserID),
		}
		if t := idx.match(v.EditionID, v.TicketTypeID, v.Origin); t != nil {
			id := t.ID
			name := t.Name
			item.TemplateID = &id
			item.TemplateName = &name
			item.DesignData = t.DesignData
		} else {
			item.DesignData = json.RawMessage(`{}`)
		}
		items = append(items, item)
	}
	return items, nil
}

// profileURL is the attendee's public profile on the Univents front: what a
// badge's QR points at.
func (o *Operations) profileURL(userID uuid.UUID) string {
	return strings.TrimRight(o.AppURL, "/") + "/profile/" + userID.String()
}

// RevokeForRegistration revokes a participant badge when its registration is
// cancelled or expires. Called by the checkout feature.
func (o *Operations) RevokeForRegistration(ctx context.Context, registrationID uuid.UUID, reason string) error {
	return o.Emissions.RevokeByRegistration(ctx, registrationID, reason)
}

// RevokeStaffBadgesForUser revokes the user's staff badges on current and
// future editions of the event. Past-edition badges are kept as keepsakes.
func (o *Operations) RevokeStaffBadgesForUser(ctx context.Context, eventID, userID uuid.UUID) error {
	editions, err := o.Editions.ListPublic(ctx, eventID)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, edition := range editions {
		if edition.EndsAt.Before(now) {
			continue
		}
		err = o.Emissions.Revoke(ctx, edition.ID, userID, models.BadgeEmissionOriginStaff, "staff removed")
		if err != nil {
			return err
		}
	}
	return nil
}

// editionTemplates indexes one edition's templates for derivation:
// a specific ticket-type template wins over the edition default (NULL
// ticket_type_id); the staff-only template (origin 'staff') wins for
// staff-origin emissions; no match at all yields no template (placeholder
// badge).
type editionTemplates struct {
	defaultT *models.BadgeTemplate
	staffT   *models.BadgeTemplate
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
		switch {
		case t.Origin != nil && *t.Origin == models.BadgeTemplateOriginStaff:
			et.staffT = t
		case t.TicketTypeID == nil:
			et.defaultT = t
		default:
			et.byTicket[*t.TicketTypeID] = t
		}
		idx[t.EditionID] = et
	}
	return idx
}

func (idx templateIndex) match(editionID uuid.UUID, ticketTypeID *uuid.UUID, origin models.BadgeEmissionOrigin) *models.BadgeTemplate {
	et, ok := idx[editionID]
	if !ok {
		return nil
	}
	if origin == models.BadgeEmissionOriginStaff && et.staffT != nil {
		return et.staffT
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
	templates, err := o.Repo.ListByEditionIDs(ctx, editionIDs)
	if err != nil {
		return nil, err
	}
	return indexTemplates(templates), nil
}

func (o *Operations) profileBadge(v models.BadgeEmissionView, idx templateIndex) models.BadgeProfileBadge {
	t := idx.match(v.EditionID, v.TicketTypeID, v.Origin)
	badge := models.BadgeProfileBadge{
		EmissionID:  v.ID,
		EditionID:   v.EditionID,
		EditionName: v.EditionName,
		EventName:   v.EventName,
		Origin:      v.Origin,
		TicketName:  v.TicketName,
		ActionURL:   o.profileURL(v.UserID),
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

func (o *Operations) UpdateTemplate(ctx context.Context, input models.UpdateBadgeTemplateInput) (*models.BadgeTemplate, error) {
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
