package badges_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"

	"univents/models"
	"univents/ports"
)

func TestPrintByEdition_ActiveOnly(t *testing.T) {
	mock.SetUp(t)

	editionID := uuid.New()
	activeID := uuid.New()
	revokedID := uuid.New()

	emissions := mock.Mock[ports.BadgeEmissionRepo]()
	templates := mock.Mock[ports.BadgeTemplateRepo]()
	editions := mock.Mock[ports.EditionRepo]()
	authzEvents := mock.Mock[ports.EventRepo]()

	mock.When(editions.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(&models.Edition{ID: editionID, EventID: uuid.New()}, nil)
	mock.When(authzEvents.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
		ThenReturn(models.EventMemberRoleAdmin, nil)

	reason := "cancelled"
	mock.When(emissions.ListViewsByEdition(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn([]models.BadgeEmissionView{
			view(activeID, editionID, uuid.New(), models.BadgeEmissionOriginParticipant, "Ed", time.Now().Add(time.Hour), "Ev", nil, new("Vip")),
			{
				BadgeEmission: models.BadgeEmission{
					ID: revokedID, EditionID: editionID, UserID: uuid.New(),
					Origin: models.BadgeEmissionOriginParticipant,
					Status: models.BadgeEmissionStatusRevoked, StatusReason: &reason,
					EmittedAt: time.Now(),
				},
				EditionName: "Ed", EndsAt: time.Now().Add(time.Hour), EventName: "Ev",
			},
		}, nil)

	tpl := models.BadgeTemplate{ID: uuid.New(), EditionID: editionID, Name: "Default", DesignData: []byte(`{"bg":"#fff"}`)}
	mock.When(templates.ListByEditionIDs(mock.AnyContext(), mock.Any[[]uuid.UUID]())).
		ThenReturn([]models.BadgeTemplate{tpl}, nil)

	ops := newOpsWithAuthz(t, templates, emissions, editions, authzEvents)

	items, err := ops.PrintByEdition(ownerCtx(), editionID, nil)
	if err != nil {
		t.Fatalf("PrintByEdition: %v", err)
	}

	if len(items) != 1 || items[0].EmissionID != activeID {
		t.Fatalf("revoked must not print; want only %s, got %+v", activeID, items)
	}
	if items[0].TemplateID == nil || *items[0].TemplateID != tpl.ID {
		t.Errorf("want derived template, got %v", items[0].TemplateID)
	}
	if string(items[0].DesignData) != `{"bg":"#fff"}` {
		t.Errorf("want design data in payload, got %s", items[0].DesignData)
	}
	if items[0].ActionURL == "" || items[0].EditionName != "Ed" || items[0].EventName != "Ev" {
		t.Errorf("want resolved variables, got %+v", items[0])
	}
}

func TestPrintByEdition_FiltersByEmissionIDs(t *testing.T) {
	mock.SetUp(t)

	editionID := uuid.New()
	firstID := uuid.New()
	secondID := uuid.New()

	emissions := mock.Mock[ports.BadgeEmissionRepo]()
	templates := mock.Mock[ports.BadgeTemplateRepo]()
	editions := mock.Mock[ports.EditionRepo]()
	authzEvents := mock.Mock[ports.EventRepo]()

	mock.When(editions.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(&models.Edition{ID: editionID, EventID: uuid.New()}, nil)
	mock.When(authzEvents.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
		ThenReturn(models.EventMemberRoleOwner, nil)

	mock.When(emissions.ListViewsByEdition(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn([]models.BadgeEmissionView{
			view(firstID, editionID, uuid.New(), models.BadgeEmissionOriginStaff, "Ed", time.Now().Add(time.Hour), "Ev", nil, nil),
			view(secondID, editionID, uuid.New(), models.BadgeEmissionOriginStaff, "Ed", time.Now().Add(time.Hour), "Ev", nil, nil),
		}, nil)
	mock.When(templates.ListByEditionIDs(mock.AnyContext(), mock.Any[[]uuid.UUID]())).
		ThenReturn(nil, nil)

	ops := newOpsWithAuthz(t, templates, emissions, editions, authzEvents)

	items, err := ops.PrintByEdition(ownerCtx(), editionID, []uuid.UUID{secondID})
	if err != nil {
		t.Fatalf("PrintByEdition: %v", err)
	}

	if len(items) != 1 || items[0].EmissionID != secondID {
		t.Fatalf("want only emission %s, got %+v", secondID, items)
	}
	// No template matched: placeholder design for the renderer.
	if string(items[0].DesignData) != "{}" {
		t.Errorf("want placeholder design {}, got %s", items[0].DesignData)
	}
}
