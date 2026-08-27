package badges_test

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"

	"univents/models"
	"univents/ports"
)

func TestUpdateTemplate_MergesPartialFields(t *testing.T) {
	cases := []struct {
		name       string
		bodyName   *string
		bodyDesign json.RawMessage
		wantName   string
		wantDesign string
	}{
		{"name only", new("New Name"), nil, "New Name", `{"old":1}`},
		{"design only", nil, json.RawMessage(`{"new":2}`), "Old Name", `{"new":2}`},
		{"both", new("New Name"), json.RawMessage(`{"new":2}`), "New Name", `{"new":2}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mock.SetUp(t)

			templateID := uuid.New()
			editionID := uuid.New()
			oldTemplate := &models.BadgeTemplate{
				ID: templateID, EditionID: editionID,
				Name: "Old Name", DesignData: json.RawMessage(`{"old":1}`),
			}

			templates := mock.Mock[ports.BadgeTemplateRepo]()
			emissions := mock.Mock[ports.BadgeEmissionRepo]()
			editions := mock.Mock[ports.EditionRepo]()
			authzEvents := mock.Mock[ports.EventRepo]()

			mock.When(templates.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(oldTemplate, nil)
			mock.When(editions.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
				ThenReturn(&models.Edition{ID: editionID, EventID: uuid.New()}, nil)
			mock.When(authzEvents.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
				ThenReturn(models.EventMemberRoleAdmin, nil)

			var gotID uuid.UUID
			var gotName string
			var gotDesign json.RawMessage
			mock.When(templates.Update(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[string](), mock.Any[json.RawMessage]())).
				ThenAnswer(func(args []any) []any {
					gotID = args[1].(uuid.UUID)
					gotName = args[2].(string)
					gotDesign = args[3].(json.RawMessage)
					return []any{&models.BadgeTemplate{ID: templateID, Name: gotName, DesignData: gotDesign}, nil}
				})

			ops := newOpsWithAuthz(t, templates, emissions, editions, authzEvents)

			_, err := ops.UpdateTemplate(ownerCtx(), models.UpdateBadgeTemplateInput{
				TemplateID: templateID,
				Name:       tc.bodyName,
				DesignData: tc.bodyDesign,
			})
			if err != nil {
				t.Fatalf("UpdateTemplate: %v", err)
			}

			if gotID != templateID {
				t.Errorf("want template id %s, got %s", templateID, gotID)
			}
			if gotName != tc.wantName {
				t.Errorf("want name %q, got %q", tc.wantName, gotName)
			}
			if string(gotDesign) != tc.wantDesign {
				t.Errorf("want design %s, got %s", tc.wantDesign, gotDesign)
			}
		})
	}
}

func TestUpdateTemplate_EmptyPatchIsNoOpMerge(t *testing.T) {
	mock.SetUp(t)

	templateID := uuid.New()
	oldTemplate := &models.BadgeTemplate{
		ID: templateID, EditionID: uuid.New(),
		Name: "Old Name", DesignData: json.RawMessage(`{"old":1}`),
	}

	templates := mock.Mock[ports.BadgeTemplateRepo]()
	emissions := mock.Mock[ports.BadgeEmissionRepo]()
	editions := mock.Mock[ports.EditionRepo]()
	authzEvents := mock.Mock[ports.EventRepo]()

	mock.When(templates.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).ThenReturn(oldTemplate, nil)
	mock.When(editions.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenReturn(&models.Edition{ID: oldTemplate.EditionID, EventID: uuid.New()}, nil)
	mock.When(authzEvents.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
		ThenReturn(models.EventMemberRoleOwner, nil)

	var gotName string
	var gotDesign json.RawMessage
	mock.When(templates.Update(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[string](), mock.Any[json.RawMessage]())).
		ThenAnswer(func(args []any) []any {
			gotName = args[2].(string)
			gotDesign = args[3].(json.RawMessage)
			return []any{&models.BadgeTemplate{ID: templateID, Name: gotName, DesignData: gotDesign}, nil}
		})

	ops := newOpsWithAuthz(t, templates, emissions, editions, authzEvents)

	_, err := ops.UpdateTemplate(ownerCtx(), models.UpdateBadgeTemplateInput{TemplateID: templateID})
	if err != nil {
		t.Fatalf("UpdateTemplate: %v", err)
	}

	if gotName != "Old Name" || string(gotDesign) != `{"old":1}` {
		t.Errorf("empty PATCH must keep values, got %q %s", gotName, gotDesign)
	}
}
