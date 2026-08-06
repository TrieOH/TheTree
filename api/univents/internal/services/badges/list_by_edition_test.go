package badges_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"

	"univents/models"
	"univents/ports"
)

func TestListByEdition_IncludesAllStatusesAndDerivesTemplate(t *testing.T) {
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

	reason := "staff removed"
	mock.When(emissions.ListViewsByEdition(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenAnswer(func(_ []any) []any {
			return []any{[]models.BadgeEmissionView{
				view(activeID, editionID, uuid.New(), models.BadgeEmissionOriginParticipant, "Ed", time.Now().Add(time.Hour), "Ev", nil, new("Vip")),
				{
					BadgeEmission: models.BadgeEmission{
						ID: revokedID, EditionID: editionID, UserID: uuid.New(),
						Origin: models.BadgeEmissionOriginStaff,
						Status: models.BadgeEmissionStatusRevoked, StatusReason: &reason,
						EmittedAt: time.Now(),
					},
					EditionName: "Ed", EndsAt: time.Now().Add(time.Hour), EventName: "Ev",
				},
			}, nil}
		})

	tpl := models.BadgeTemplate{ID: uuid.New(), EditionID: editionID, Name: "Default", DesignData: []byte(`{}`)}
	mock.When(templates.ListByEditionIDs(mock.AnyContext(), mock.Any[[]uuid.UUID]())).
		ThenReturn([]models.BadgeTemplate{tpl}, nil)

	ops := newOpsWithAuthz(t, templates, emissions, nil, editions, nil, nil, authzEvents)

	items, err := ops.ListByEdition(ownerCtx(), editionID)
	if err != nil {
		t.Fatalf("ListByEdition: %v", err)
	}

	if len(items) != 2 {
		t.Fatalf("want 2 items (active + revoked), got %d", len(items))
	}
	byID := map[uuid.UUID]models.BadgeEditionEmission{}
	for _, it := range items {
		byID[it.ID] = it
	}
	if byID[revokedID].Status != models.BadgeEmissionStatusRevoked {
		t.Error("revoked emission must be included in browse, with revoked status")
	}
	if byID[activeID].TemplateID == nil || *byID[activeID].TemplateID != tpl.ID {
		t.Errorf("want derived template on active item, got %v", byID[activeID].TemplateID)
	}
}
