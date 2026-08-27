package badges_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"

	"univents/models"
	"univents/ports"
)

func staffTemplate(id, editionID uuid.UUID, name string) models.BadgeTemplate {
	origin := models.BadgeTemplateOriginStaff
	return models.BadgeTemplate{
		ID: id, EditionID: editionID, Origin: &origin,
		Name: name, DesignData: json.RawMessage(`{"kind":"staff"}`),
	}
}

func TestListByUser_StaffTemplateWinsForStaff(t *testing.T) {
	mock.SetUp(t)

	editionID := uuid.New()
	userID := uuid.New()
	staffTpl := staffTemplate(uuid.New(), editionID, "StaffTpl")
	defaultTpl := models.BadgeTemplate{ID: uuid.New(), EditionID: editionID, Name: "Default", DesignData: json.RawMessage(`{}`)}

	emissions := mock.Mock[ports.BadgeEmissionRepo]()
	templates := mock.Mock[ports.BadgeTemplateRepo]()

	mock.When(emissions.ListViewsByUser(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn([]models.BadgeEmissionView{
			view(uuid.New(), editionID, userID, models.BadgeEmissionOriginStaff, "Ed", time.Now().Add(time.Hour), "Ev", nil, nil),
		}, nil)
	mock.When(templates.ListByEditionIDs(mock.AnyContext(), mock.Any[[]uuid.UUID]())).
		ThenReturn([]models.BadgeTemplate{defaultTpl, staffTpl}, nil)

	ops := newOps(t, templates, emissions, nil, nil, nil, nil)

	groups, err := ops.ListByUser(context.Background(), userID)
	if err != nil {
		t.Fatalf("ListByUser: %v", err)
	}

	got := groups.Staff.Current[0]
	if got.TemplateID == nil || *got.TemplateID != staffTpl.ID {
		t.Errorf("want staff template %s for staff emission, got %v", staffTpl.ID, got.TemplateID)
	}
	if got.TemplateName == nil || *got.TemplateName != "StaffTpl" {
		t.Errorf("want staff template name, got %v", got.TemplateName)
	}
}

func TestListByUser_StaffFallsBackToDefault(t *testing.T) {
	mock.SetUp(t)

	editionID := uuid.New()
	userID := uuid.New()
	defaultTpl := models.BadgeTemplate{ID: uuid.New(), EditionID: editionID, Name: "Default", DesignData: json.RawMessage(`{}`)}

	emissions := mock.Mock[ports.BadgeEmissionRepo]()
	templates := mock.Mock[ports.BadgeTemplateRepo]()

	mock.When(emissions.ListViewsByUser(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn([]models.BadgeEmissionView{
			view(uuid.New(), editionID, userID, models.BadgeEmissionOriginStaff, "Ed", time.Now().Add(time.Hour), "Ev", nil, nil),
		}, nil)
	mock.When(templates.ListByEditionIDs(mock.AnyContext(), mock.Any[[]uuid.UUID]())).
		ThenReturn([]models.BadgeTemplate{defaultTpl}, nil)

	ops := newOps(t, templates, emissions, nil, nil, nil, nil)

	groups, err := ops.ListByUser(context.Background(), userID)
	if err != nil {
		t.Fatalf("ListByUser: %v", err)
	}

	// No staff template: staff keeps the pre-feature behavior — edition default.
	got := groups.Staff.Current[0]
	if got.TemplateID == nil || *got.TemplateID != defaultTpl.ID {
		t.Errorf("want default template for staff emission, got %v", got.TemplateID)
	}
}

func TestListByUser_StaffTemplateDoesNotLeakToAttendant(t *testing.T) {
	mock.SetUp(t)

	editionID := uuid.New()
	ticketTypeID := uuid.New()
	userID := uuid.New()
	staffTpl := staffTemplate(uuid.New(), editionID, "StaffTpl")
	defaultTpl := models.BadgeTemplate{ID: uuid.New(), EditionID: editionID, Name: "Default", DesignData: json.RawMessage(`{}`)}
	specificTpl := models.BadgeTemplate{ID: uuid.New(), EditionID: editionID, TicketTypeID: &ticketTypeID, Name: "VipTpl", DesignData: json.RawMessage(`{}`)}

	emissions := mock.Mock[ports.BadgeEmissionRepo]()
	templates := mock.Mock[ports.BadgeTemplateRepo]()

	mock.When(emissions.ListViewsByUser(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn([]models.BadgeEmissionView{
			view(uuid.New(), editionID, userID, models.BadgeEmissionOriginParticipant, "Ed", time.Now().Add(time.Hour), "Ev", &ticketTypeID, new("Vip")),
			view(uuid.New(), editionID, userID, models.BadgeEmissionOriginParticipant, "Ed", time.Now().Add(time.Hour), "Ev", nil, nil),
		}, nil)
	mock.When(templates.ListByEditionIDs(mock.AnyContext(), mock.Any[[]uuid.UUID]())).
		ThenReturn([]models.BadgeTemplate{defaultTpl, staffTpl, specificTpl}, nil)

	ops := newOps(t, templates, emissions, nil, nil, nil, nil)

	groups, err := ops.ListByUser(context.Background(), userID)
	if err != nil {
		t.Fatalf("ListByUser: %v", err)
	}

	// The staff template is never offered to attendant emissions.
	if got := groups.Attendant.Current[0]; got.TemplateID == nil || *got.TemplateID != specificTpl.ID {
		t.Errorf("want ticket-specific template for attendant, got %v", got.TemplateID)
	}
	if got := groups.Attendant.Current[1]; got.TemplateID == nil || *got.TemplateID != defaultTpl.ID {
		t.Errorf("want default template for ticket-less attendant, got %v", got.TemplateID)
	}
}

func TestPrintByEdition_StaffTemplateDesignData(t *testing.T) {
	mock.SetUp(t)

	editionID := uuid.New()
	staffTpl := staffTemplate(uuid.New(), editionID, "StaffTpl")

	emissions := mock.Mock[ports.BadgeEmissionRepo]()
	templates := mock.Mock[ports.BadgeTemplateRepo]()
	editions := mock.Mock[ports.EditionRepo]()
	authzEvents := mock.Mock[ports.EventRepo]()

	mock.When(editions.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(&models.Edition{ID: editionID, EventID: uuid.New()}, nil)
	mock.When(authzEvents.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
		ThenReturn(models.EventMemberRoleAdmin, nil)
	mock.When(emissions.ListViewsByEdition(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn([]models.BadgeEmissionView{
			view(uuid.New(), editionID, uuid.New(), models.BadgeEmissionOriginStaff, "Ed", time.Now().Add(time.Hour), "Ev", nil, nil),
		}, nil)
	mock.When(templates.ListByEditionIDs(mock.AnyContext(), mock.Any[[]uuid.UUID]())).
		ThenReturn([]models.BadgeTemplate{staffTpl}, nil)

	ops := newOpsWithAuthz(t, templates, emissions, editions, authzEvents)

	items, err := ops.PrintByEdition(ownerCtx(), editionID, nil)
	if err != nil {
		t.Fatalf("PrintByEdition: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("want 1 print item, got %d", len(items))
	}
	if items[0].TemplateID == nil || *items[0].TemplateID != staffTpl.ID {
		t.Errorf("want staff template in print payload, got %v", items[0].TemplateID)
	}
	if string(items[0].DesignData) != `{"kind":"staff"}` {
		t.Errorf("want staff design data, got %s", items[0].DesignData)
	}
}

func TestCreateTemplate_RejectsInvalidStaffScope(t *testing.T) {
	mock.SetUp(t)

	// Validation runs before any repo/authz access, so no mocks are needed.
	ops := newOps(t, nil, nil, nil, nil, nil, nil)

	staffOrigin := models.BadgeTemplateOriginStaff

	cases := []struct {
		name       string
		origin     *models.BadgeTemplateOrigin
		ticketType *uuid.UUID
	}{
		{"staff with ticket type", &staffOrigin, new(uuid.UUID)},
		{"unknown origin", new(models.BadgeTemplateOrigin), nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ops.CreateTemplate(context.Background(), models.CreateBadgeTemplateInput{
				EditionID:    uuid.New(),
				TicketTypeID: tc.ticketType,
				Origin:       tc.origin,
				Name:         "Staff",
				DesignData:   json.RawMessage(`{}`),
			})
			if err == nil {
				t.Fatalf("want error for scope %v, got nil", tc.origin)
			}
		})
	}
}

func TestCreateTemplate_AcceptsStaffScope(t *testing.T) {
	mock.SetUp(t)

	editionID := uuid.New()
	staffOrigin := models.BadgeTemplateOriginStaff

	templates := mock.Mock[ports.BadgeTemplateRepo]()
	emissions := mock.Mock[ports.BadgeEmissionRepo]()
	editions := mock.Mock[ports.EditionRepo]()
	authzEvents := mock.Mock[ports.EventRepo]()

	mock.When(editions.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(&models.Edition{ID: editionID, EventID: uuid.New()}, nil)
	mock.When(authzEvents.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
		ThenReturn(models.EventMemberRoleOwner, nil)

	var gotOrigin *models.BadgeTemplateOrigin
	mock.When(templates.Create(mock.AnyContext(), mock.Any[*models.BadgeTemplate]())).
		ThenAnswer(func(args []any) []any {
			in := args[1].(*models.BadgeTemplate)
			gotOrigin = in.Origin
			return []any{in, nil}
		})

	ops := newOpsWithAuthz(t, templates, emissions, editions, authzEvents)

	_, err := ops.CreateTemplate(ownerCtx(), models.CreateBadgeTemplateInput{
		EditionID:  editionID,
		Origin:     &staffOrigin,
		Name:       "Staff",
		DesignData: json.RawMessage(`{}`),
	})
	if err != nil {
		t.Fatalf("CreateTemplate: %v", err)
	}

	if gotOrigin == nil || *gotOrigin != models.BadgeTemplateOriginStaff {
		t.Errorf("want staff origin passed to repo, got %v", gotOrigin)
	}
}
