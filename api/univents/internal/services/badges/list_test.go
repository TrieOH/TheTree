package badges_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"

	"univents/models"
	"univents/ports"
)

func TestListByUser_GroupsAndSorts(t *testing.T) {
	mock.SetUp(t)

	now := time.Now()
	editionCurrent := uuid.New()
	editionSoon := uuid.New()
	editionPast := uuid.New()
	userID := uuid.New()

	emissions := mock.Mock[ports.BadgeEmissionRepo]()
	templates := mock.Mock[ports.BadgeTemplateRepo]()

	mock.When(emissions.ListViewsByUser(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn([]models.BadgeEmissionView{
			view(uuid.New(), editionCurrent, userID, models.BadgeEmissionOriginParticipant, "Current", now.Add(48*time.Hour), "Event", nil, new("Vip")),
			view(uuid.New(), editionSoon, userID, models.BadgeEmissionOriginParticipant, "Soon", now.Add(24*time.Hour), "Event", nil, nil),
			view(uuid.New(), editionPast, userID, models.BadgeEmissionOriginParticipant, "Past", now.Add(-48*time.Hour), "Event", nil, nil),
			view(uuid.New(), editionCurrent, userID, models.BadgeEmissionOriginStaff, "Current", now.Add(48*time.Hour), "Event", nil, nil),
		}, nil)

	mock.When(templates.ListByEditionIDs(mock.AnyContext(), mock.Any[[]uuid.UUID]())).
		ThenReturn([]models.BadgeTemplate{
			{ID: uuid.New(), EditionID: editionCurrent, Name: "Default", DesignData: []byte(`{"a":1}`)},
			{ID: uuid.New(), EditionID: editionSoon, Name: "SoonTpl", DesignData: []byte(`{"b":2}`)},
		}, nil)

	ops := newOps(t, templates, emissions, nil, nil, nil, nil)

	groups, err := ops.ListByUser(context.Background(), userID)
	if err != nil {
		t.Fatalf("ListByUser: %v", err)
	}

	if len(groups.Attendant.Current) != 2 {
		t.Fatalf("want 2 current attendant badges, got %d", len(groups.Attendant.Current))
	}
	// Current sorts soonest-ending first ("most current").
	if groups.Attendant.Current[0].EditionName != "Soon" {
		t.Errorf("most current should be 'Soon', got %q", groups.Attendant.Current[0].EditionName)
	}
	if len(groups.Attendant.Past) != 1 || groups.Attendant.Past[0].EditionName != "Past" {
		t.Errorf("want past badge 'Past', got %+v", groups.Attendant.Past)
	}
	if len(groups.Staff.Current) != 1 || groups.Staff.Current[0].EditionName != "Current" {
		t.Errorf("want staff current badge, got %+v", groups.Staff.Current)
	}

	// Template derivation: specific ticket template wins over the edition default.
	if groups.Attendant.Current[1].TemplateName == nil || *groups.Attendant.Current[1].TemplateName != "Default" {
		t.Errorf("want derived default template, got %+v", groups.Attendant.Current[1].TemplateName)
	}
}

func TestListByUser_DerivesSpecificOverDefault(t *testing.T) {
	mock.SetUp(t)

	editionID := uuid.New()
	ticketTypeID := uuid.New()
	userID := uuid.New()

	emissions := mock.Mock[ports.BadgeEmissionRepo]()
	templates := mock.Mock[ports.BadgeTemplateRepo]()

	mock.When(emissions.ListViewsByUser(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn([]models.BadgeEmissionView{
			view(uuid.New(), editionID, userID, models.BadgeEmissionOriginParticipant, "Ed", time.Now().Add(time.Hour), "Ev", &ticketTypeID, new("Vip")),
		}, nil)

	defaultTpl := models.BadgeTemplate{ID: uuid.New(), EditionID: editionID, Name: "Default", DesignData: []byte(`{}`)}
	specificTpl := models.BadgeTemplate{ID: uuid.New(), EditionID: editionID, Name: "VipTpl", DesignData: []byte(`{}`)}
	mock.When(templates.ListByEditionIDs(mock.AnyContext(), mock.Any[[]uuid.UUID]())).
		ThenReturn([]models.BadgeTemplate{defaultTpl, specificTpl}, nil)

	ops := newOps(t, templates, emissions, nil, nil, nil, nil)

	groups, err := ops.ListByUser(context.Background(), userID)
	if err != nil {
		t.Fatalf("ListByUser: %v", err)
	}

	got := groups.Attendant.Current[0]
	if got.TemplateID == nil || *got.TemplateID != specificTpl.ID {
		t.Errorf("want specific template %s, got %v", specificTpl.ID, got.TemplateID)
	}
	if got.TicketName == nil || *got.TicketName != "Vip" {
		t.Errorf("want ticket name Vip, got %v", got.TicketName)
	}
}

func TestListByUser_NoTemplateMeansPlaceholder(t *testing.T) {
	mock.SetUp(t)

	editionID := uuid.New()
	userID := uuid.New()

	emissions := mock.Mock[ports.BadgeEmissionRepo]()
	templates := mock.Mock[ports.BadgeTemplateRepo]()

	mock.When(emissions.ListViewsByUser(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn([]models.BadgeEmissionView{
			view(uuid.New(), editionID, userID, models.BadgeEmissionOriginStaff, "Ed", time.Now().Add(time.Hour), "Ev", nil, nil),
		}, nil)
	mock.When(templates.ListByEditionIDs(mock.AnyContext(), mock.Any[[]uuid.UUID]())).
		ThenReturn(nil, nil)

	ops := newOps(t, templates, emissions, nil, nil, nil, nil)

	groups, err := ops.ListByUser(context.Background(), userID)
	if err != nil {
		t.Fatalf("ListByUser: %v", err)
	}

	if len(groups.Staff.Current) != 1 {
		t.Fatalf("want 1 staff badge, got %d", len(groups.Staff.Current))
	}
	if groups.Staff.Current[0].TemplateID != nil {
		t.Errorf("want nil template (placeholder), got %v", groups.Staff.Current[0].TemplateID)
	}
}
