package emails

import (
	"context"
	"testing"
	"time"

	"IdentityX/internal/tokens"
	"IdentityX/models"
	"IdentityX/ports"
	"lib/crypto"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
	"github.com/ovechkin-dm/mockio/mock"
	"github.com/riverqueue/river"
)

const testHMAC = "test-hmac"

func stubSender(t *testing.T, actionTokens ports.ActionTokenRepo) (*Sender, *[]SendAuthEmailArgs) {
	t.Helper()
	var enqueued []SendAuthEmailArgs
	enqueuer := mock.Mock[Enqueuer]()
	_ = mock.When(enqueuer.Insert(mock.AnyContext(), mock.Any[river.JobArgs](), mock.Any[*river.InsertOpts]())).
		ThenAnswer(func(args []any) []any {
			enqueued = append(enqueued, args[1].(SendAuthEmailArgs))
			return []any{nil, nil}
		})
	mgr := tokens.NewActionTokenManager(actionTokens, []byte(testHMAC), tokens.ActionTokenConfig{
		VerifyTTL: 10 * time.Minute,
		ResetTTL:  10 * time.Minute,
	})
	return NewSender(mgr, "https://identityx.example.com", "IdentityX", enqueuer), &enqueued
}

func actorWithEmail() *models.Actor {
	email := "user@example.com"
	return &models.Actor{ID: uuid.New(), Email: &email, Type: models.HumanActorType}
}

func TestSenderSendVerifyPlatformActor(t *testing.T) {
	mock.SetUp(t)
	actor := actorWithEmail()

	actionTokens := mock.Mock[ports.ActionTokenRepo]()
	_ = mock.When(actionTokens.Insert(mock.AnyContext(), mock.Any[models.ActionToken]())).
		ThenAnswer(func(args []any) []any {
			token := args[1].(models.ActionToken)
			return []any{&token, nil}
		})

	sender, enqueued := stubSender(t, actionTokens)
	err := sender.SendVerify(context.Background(), actor, nil)
	if err != nil {
		t.Fatalf("SendVerify: %v", err)
	}

	if len(*enqueued) != 1 {
		t.Fatalf("enqueued %d jobs, want 1", len(*enqueued))
	}
	job := (*enqueued)[0]
	if job.TemplateKind != string(models.VerifyEmailTemplateKind) {
		t.Fatalf("kind = %q, want verify", job.TemplateKind)
	}
	if job.BaseDomain != "https://identityx.example.com" {
		t.Fatalf("base = %q", job.BaseDomain)
	}
	if job.ProjectID != nil {
		t.Fatalf("platform actor must not carry project_id")
	}
	if job.Expiry != 10 {
		t.Fatalf("expiry = %d, want 10", job.Expiry)
	}

	// the token must be a verifiable HMAC JWT with the right claims
	claims := &models.ActionTokenClaims{}
	_, err = crypto.ParseHMACJWT(job.Token, claims, []byte(testHMAC))
	if err != nil {
		t.Fatalf("token does not parse: %v", err)
	}
	if claims.Subject != actor.ID.String() {
		t.Fatalf("sub = %q, want %q", claims.Subject, actor.ID)
	}
	if claims.Purpose != string(models.EmailVerifyActionTokenPurpose) {
		t.Fatalf("purpose = %q", claims.Purpose)
	}
}

func TestSenderSendResetRequiresDomain(t *testing.T) {
	mock.SetUp(t)
	actor := actorWithEmail()
	project := &models.Project{ID: uuid.New(), Name: "Acme"} // no domain

	sender, _ := stubSender(t, mock.Mock[ports.ActionTokenRepo]())
	err := sender.SendReset(context.Background(), actor, project)
	if err == nil || !fun.Is(err, fun.CodeInternal) {
		t.Fatalf("project without domain must fail, got %v", err)
	}
}

func TestSenderSendVerifyProjectActorUsesDomain(t *testing.T) {
	mock.SetUp(t)
	actor := actorWithEmail()
	domain := "acme.example.com"
	project := &models.Project{ID: uuid.New(), Name: "Acme", Domain: &domain}

	actionTokens := mock.Mock[ports.ActionTokenRepo]()
	_ = mock.When(actionTokens.Insert(mock.AnyContext(), mock.Any[models.ActionToken]())).
		ThenAnswer(func(args []any) []any {
			token := args[1].(models.ActionToken)
			return []any{&token, nil}
		})

	sender, enqueued := stubSender(t, actionTokens)
	err := sender.SendVerify(context.Background(), actor, project)
	if err != nil {
		t.Fatalf("SendVerify: %v", err)
	}
	job := (*enqueued)[0]
	if job.BaseDomain != "https://acme.example.com" {
		t.Fatalf("base = %q, want https://acme.example.com", job.BaseDomain)
	}
	if job.ProjectName != "Acme" {
		t.Fatalf("project name = %q", job.ProjectName)
	}
	if job.ProjectID == nil || *job.ProjectID != project.ID {
		t.Fatalf("project_id not carried")
	}
}
