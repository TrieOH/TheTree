package authn

import (
	"context"
	"testing"

	"IdentityX/internal/authz"
	"IdentityX/internal/emails"
	"IdentityX/internal/services/tos"
	"IdentityX/internal/tokens"
	"IdentityX/models"
	"IdentityX/ports"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"
)

// stubRegisterTos builds authn ops whose tos service reports the given
// current terms, and returns the mocks the test asserts against.
func stubRegisterTos(t *testing.T, projectID uuid.UUID, current *models.TermsOfService) (*Operations, ports.TosRepo) {
	t.Helper()

	projects := mock.Mock[ports.ProjectRepo]()
	_ = mock.When(projects.GetByID(mock.AnyContext(), mock.Equal(projectID))).
		ThenReturn(&models.Project{ID: projectID}, nil)

	actors := mock.Mock[ports.ActorRepo]()
	_ = mock.When(actors.Register(mock.AnyContext(), mock.Any[models.Actor]())).
		ThenAnswer(func(args []any) []any {
			a := args[1].(models.Actor)
			return []any{&models.Actor{ID: uuid.New(), Email: a.Email, Type: a.Type}, nil}
		})

	tosRepo := mock.Mock[ports.TosRepo]()
	_ = mock.When(tosRepo.GetCurrent(mock.AnyContext(), mock.Equal(projectID))).
		ThenReturn(current, nil)

	sender := mock.Mock[ports.EmailSender]()
	_ = mock.When(sender.SendVerify(mock.AnyContext(), mock.Any[*models.Actor](), mock.Any[*models.Project]())).
		ThenReturn(nil)

	tosOps := tos.NewOperations(tosRepo, projects, actors,
		authz.New(mock.Mock[ports.OrganizationRepo](), projects, mock.Mock[ports.PlatformRolesRepo]()),
		emails.NewTosNotifier(mock.Mock[emails.Enqueuer]()),
	)
	ops := NewOperations(
		actors,
		projects,
		mock.Mock[ports.PlatformRolesRepo](),
		tokens.NewManager(mock.Mock[ports.CryptoKeysRepo](), mock.Mock[ports.BlacklistRepo](), mock.Mock[ports.ActorRepo](), mock.Mock[ports.ProjectRepo](), tokens.Config{}),
		actionMgr(mock.Mock[ports.ActionTokenRepo]()),
		sender,
		tosOps,
	)
	return ops, tosRepo
}

func TestRegisterRejectsMissingTosAcceptance(t *testing.T) {
	mock.SetUp(t)
	projectID := uuid.New()
	current := &models.TermsOfService{ID: uuid.New(), ProjectID: projectID, Version: 1, Content: "termos v1"}
	ops, _ := stubRegisterTos(t, projectID, current)
	err := ops.Register(context.Background(), models.IDXRegisterInput{
		Email:     "user@acme.com.br",
		Password:  "S3curePassw0rd!",
		ProjectID: &projectID,
	})
	if err == nil || !fun.Is(err, fun.CodeValidation) {
		t.Fatalf("registration into a project with terms without accepted_tos must fail, got %v", err)
	}
}

func TestRegisterRecordsAcceptanceWhenAccepted(t *testing.T) {
	mock.SetUp(t)
	projectID := uuid.New()
	current := &models.TermsOfService{ID: uuid.New(), ProjectID: projectID, Version: 1, Content: "termos v1"}
	ops, tosRepo := stubRegisterTos(t, projectID, current)

	var recorded *models.TosAcceptance
	_ = mock.When(tosRepo.RecordAcceptance(mock.AnyContext(), mock.Any[models.TosAcceptance]())).
		ThenAnswer(func(args []any) []any {
			r := args[1].(models.TosAcceptance)
			recorded = &r
			return []any{recorded, nil}
		})

	err := ops.Register(context.Background(), models.IDXRegisterInput{
		Email:       "user@acme.com.br",
		Password:    "S3curePassw0rd!",
		AcceptedTos: true,
		ProjectID:   &projectID,
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if recorded == nil {
		t.Fatal("the acceptance row must be recorded on registration")
	}
	if recorded.TosVersion != 1 {
		t.Fatalf("acceptance must pin the current version, got %d", recorded.TosVersion)
	}
}

func TestRegisterSkipsTosGateForProjectWithoutTerms(t *testing.T) {
	mock.SetUp(t)
	projectID := uuid.New()
	ops, tosRepo := stubRegisterTos(t, projectID, nil)

	// v0: GetCurrent answers not-found; CurrentTos maps it to nil.
	_ = mock.When(tosRepo.GetCurrent(mock.AnyContext(), mock.Equal(projectID))).
		ThenReturn(nil, fun.ErrNotFound("terms_of_service not found"))

	err := ops.Register(context.Background(), models.IDXRegisterInput{
		Email:     "user@acme.com.br",
		Password:  "S3curePassw0rd!",
		ProjectID: &projectID,
	})
	if err != nil {
		t.Fatalf("registration without terms must not gate, got %v", err)
	}
}
