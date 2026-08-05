package email_templates

import (
	"context"
	"testing"

	"IdentityX/internal/authz"
	"IdentityX/models"
	"IdentityX/ports"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"
)

// adminCtx carries an admin identity so authorizeAdmin passes.
func adminCtx() context.Context {
	actor := models.Actor{ID: uuid.New(), Type: models.HumanActorType}
	return models.WithIdentity(context.Background(), &models.Identity{
		Sub:  models.Subject{ID: actor.ID, Type: models.HumanActorType},
		Cred: models.Credential{Type: models.TokenCredentialType},
	})
}

func stubAdminOps(t *testing.T, templates ports.EmailTemplateRepo) *Operations {
	t.Helper()
	projects := mock.Mock[ports.ProjectRepo]()
	_ = mock.When(projects.GetByID(mock.AnyContext(), mock.Any[uuid.UUID]())).
		ThenAnswer(func(args []any) []any {
			return []any{&models.Project{ID: args[1].(uuid.UUID)}, nil}
		})
	_ = mock.When(projects.GetRole(mock.AnyContext(), mock.Any[uuid.UUID](), mock.Any[uuid.UUID]())).
		ThenReturn(models.ProjectRoleAdmin, nil)

	ops := NewOperations(templates, projects, authz.New(mock.Mock[ports.OrganizationRepo](), projects))
	return ops
}

func TestUpsertRejectsTemplateWithoutActionURL(t *testing.T) {
	mock.SetUp(t)
	ops := stubAdminOps(t, mock.Mock[ports.EmailTemplateRepo]())

	_, err := ops.Upsert(adminCtx(), uuid.New(), models.VerifyEmailTemplateKind,
		"Verify your email",
		"<p>Please verify your account by clicking the link in your email client.</p>",
	)
	if err == nil || !fun.Is(err, fun.CodeValidation) {
		t.Fatalf("template without {{.ActionURL}} must be rejected, got %v", err)
	}
}

func TestUpsertRejectsBrokenTemplate(t *testing.T) {
	mock.SetUp(t)
	ops := stubAdminOps(t, mock.Mock[ports.EmailTemplateRepo]())

	_, err := ops.Upsert(adminCtx(), uuid.New(), models.VerifyEmailTemplateKind,
		"Verify your email",
		"<a href=\"{{.ActionURL\">verify</a>",
	)
	if err == nil {
		t.Fatal("unparseable template must be rejected")
	}
}

func TestUpsertAcceptsValidTemplateAndPersists(t *testing.T) {
	mock.SetUp(t)
	projectID := uuid.New()
	templates := mock.Mock[ports.EmailTemplateRepo]()
	_ = mock.When(templates.Upsert(mock.AnyContext(), mock.Any[models.EmailTemplate]())).
		ThenAnswer(func(args []any) []any {
			in := args[1].(models.EmailTemplate)
			return []any{&models.EmailTemplate{
				ID: uuid.New(), ProjectID: in.ProjectID, Kind: in.Kind,
				Subject: in.Subject, Body: in.Body,
			}, nil}
		})
	ops := stubAdminOps(t, templates)

	out, err := ops.Upsert(adminCtx(), projectID, models.VerifyEmailTemplateKind,
		"Welcome to Acme",
		"<p>Hi!</p><p><a href=\"{{.ActionURL}}\">Confirm my email</a></p>",
	)
	if err != nil {
		t.Fatalf("Upsert: %v", err)
	}
	if out.Source != "override" || out.Kind != models.VerifyEmailTemplateKind {
		t.Fatalf("unexpected effective template: %+v", out)
	}
}

func TestGetFallsBackToDefault(t *testing.T) {
	mock.SetUp(t)
	projectID := uuid.New()
	templates := mock.Mock[ports.EmailTemplateRepo]()
	_ = mock.When(templates.GetByProjectAndKind(mock.AnyContext(), mock.Equal(projectID), mock.Equal(models.ResetEmailTemplateKind))).
		ThenReturn(nil, fun.ErrNotFound("not found"))
	ops := stubAdminOps(t, templates)

	out, err := ops.Get(adminCtx(), projectID, models.ResetEmailTemplateKind)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if out.Source != "default" {
		t.Fatalf("source = %q, want default", out.Source)
	}
	if out.Body == "" {
		t.Fatal("default body is empty")
	}
}
