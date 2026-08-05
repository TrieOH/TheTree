package emails

import (
	"context"
	"fmt"
	"strings"
	"time"

	"IdentityX/models"
	"IdentityX/ports"
	"lib/crypto"

	"github.com/MintzyG/fun"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
)

// Enqueuer is the River insert seam the Sender depends on; *river.Client
// satisfies it, and tests mock it.
type Enqueuer interface {
	Insert(ctx context.Context, args river.JobArgs, opts *river.InsertOpts) (*rivertype.JobInsertResult, error)
}

// Sender mints single-use action tokens, persists them for anti-replay,
// and enqueues the async email job. Construction happens once at startup;
// the TTLs, app URL, and HMAC secret come from app config.
type Sender struct {
	actionTokens ports.ActionTokenRepo
	hmacSecret   []byte
	verifyTTL    time.Duration
	resetTTL     time.Duration
	appURL       string
	appName      string
	enqueuer     Enqueuer
}

func NewSender(
	actionTokens ports.ActionTokenRepo,
	hmacSecret []byte,
	verifyTTL, resetTTL time.Duration,
	appURL, appName string,
	enqueuer Enqueuer,
) *Sender {
	return &Sender{
		actionTokens: actionTokens,
		hmacSecret:   hmacSecret,
		verifyTTL:    verifyTTL,
		resetTTL:     resetTTL,
		appURL:       appURL,
		appName:      appName,
		enqueuer:     enqueuer,
	}
}

// SendVerify enqueues a verification email for the actor. project nil means
// a platform-level account (links point at APP_URL).
func (s *Sender) SendVerify(ctx context.Context, actor *models.Actor, project *models.Project) error {
	return s.send(ctx, models.VerifyEmailTemplateKind, models.EmailVerifyActionTokenPurpose, s.verifyTTL, actor, project)
}

// SendReset enqueues a password-reset email for the actor.
func (s *Sender) SendReset(ctx context.Context, actor *models.Actor, project *models.Project) error {
	return s.send(ctx, models.ResetEmailTemplateKind, models.PasswordResetActionTokenPurpose, s.resetTTL, actor, project)
}

func (s *Sender) send(
	ctx context.Context,
	kind models.EmailTemplateKind,
	purpose models.ActionTokenPurpose,
	ttl time.Duration,
	actor *models.Actor,
	project *models.Project,
) error {
	if actor.Email == nil {
		return fun.ErrInternal(fmt.Sprintf("cannot send %s email for actor %s: no email address", kind, actor.ID))
	}

	baseDomain, projectName, err := s.resolveBase(ctx, kind, project)
	if err != nil {
		return err
	}

	jti := uuid.New()
	expiresAt := time.Now().Add(ttl)
	jobProjectID := s.jobProjectID(actor, project)
	token, err := crypto.SignHMACJWT(models.ActionTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   actor.ID.String(),
			ID:        jti.String(),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Purpose:   string(purpose),
		ProjectID: jobProjectID,
	}, s.hmacSecret)
	if err != nil {
		return err
	}

	_, err = s.actionTokens.Insert(ctx, models.ActionToken{
		JTI:       jti,
		Purpose:   purpose,
		ActorID:   actor.ID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return err
	}

	_, err = s.enqueuer.Insert(ctx, SendAuthEmailArgs{
		TemplateKind: string(kind),
		Token:        token,
		ToEmail:      *actor.Email,
		ProjectID:    jobProjectID,
		ProjectName:  projectName,
		BaseDomain:   baseDomain,
		Expiry:       int(ttl.Minutes()),
	}, nil)
	return err
}

// jobProjectID prefers the resolved project: a project-scoped actor always
// carries project_id, but the project is the authority when present.
func (s *Sender) jobProjectID(actor *models.Actor, project *models.Project) *uuid.UUID {
	if project != nil {
		return &project.ID
	}
	return actor.ProjectID
}

// resolveBase picks the link base: the project's domain for project
// actors, APP_URL for platform actors. A project without a domain is a
// misconfigured tenant — fail loudly rather than mail broken links.
func (s *Sender) resolveBase(_ context.Context, kind models.EmailTemplateKind, project *models.Project) (baseDomain, projectName string, err error) {
	if project == nil {
		return s.appURL, s.appName, nil
	}
	if project.Domain == nil {
		return "", "", fun.ErrInternal(fmt.Sprintf(
			"cannot send %s email for project %s: no domain configured",
			kind, project.ID,
		))
	}
	domain := *project.Domain
	if !strings.Contains(domain, "://") {
		domain = "https://" + domain
	}
	return domain, project.Name, nil
}
