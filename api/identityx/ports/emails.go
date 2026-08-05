package ports

import (
	"IdentityX/models"
	"context"

	"github.com/google/uuid"
)

type EmailTemplateRepo interface {
	Upsert(ctx context.Context, template models.EmailTemplate) (*models.EmailTemplate, error)
	GetByProjectAndKind(ctx context.Context, projectID uuid.UUID, kind models.EmailTemplateKind) (*models.EmailTemplate, error)
	Delete(ctx context.Context, projectID uuid.UUID, kind models.EmailTemplateKind) error
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]models.EmailTemplate, error)
}

type ActionTokenRepo interface {
	Insert(ctx context.Context, token models.ActionToken) (*models.ActionToken, error)
	GetByJTI(ctx context.Context, jti uuid.UUID) (*models.ActionToken, error)
	// Consume atomically marks a single-use action token as used. It
	// returns the row only when the token was still unused; a concurrent
	// or repeated consumption is reported as not-found.
	Consume(ctx context.Context, jti uuid.UUID) (*models.ActionToken, error)
	DeleteExpired(ctx context.Context) error
}

// EmailSender dispatches verify/reset emails for an actor. SendVerify /
// SendReset mint the single-use action token, persist it, and enqueue the
// asynchronous email job. project is nil for platform-level actors.
type EmailSender interface {
	SendVerify(ctx context.Context, actor *models.Actor, project *models.Project) error
	SendReset(ctx context.Context, actor *models.Actor, project *models.Project) error
}
