package ports

import (
	"IdentityX/models"
	"context"
	"time"

	"github.com/google/uuid"
)

type TosRepo interface {
	// GetCurrent returns the project's current terms, or not-found when
	// the project has none (v0 — no terms, nothing to accept).
	GetCurrent(ctx context.Context, projectID uuid.UUID) (*models.TermsOfService, error)
	Create(ctx context.Context, tos models.TermsOfService) (*models.TermsOfService, error)
	// Update bumps the version and swaps the content in one statement.
	Update(ctx context.Context, projectID uuid.UUID, content string, effectiveAt time.Time) (*models.TermsOfService, error)
	// RecordAcceptance pins the actor's agreement to a version with the
	// consent source carried on the model; accepting the same version
	// again is idempotent (timestamp refreshed).
	RecordAcceptance(ctx context.Context, acceptance models.TosAcceptance) (*models.TosAcceptance, error)
	// LatestAcceptedVersion returns the newest version the actor has any
	// acceptance row for, 0 when none. The continued-use stamp compares
	// it against the current version.
	LatestAcceptedVersion(ctx context.Context, actorID, projectID uuid.UUID) (int, error)
	ListAcceptances(ctx context.Context, projectID uuid.UUID) ([]models.TosAcceptanceWithActor, error)
}

// TosNotifier enqueues the ToS-change notification fan-out. The tos
// Operations cross this seam instead of touching the queue; the River-
// backed implementation lives in internal/emails.
type TosNotifier interface {
	NotifyUpdate(ctx context.Context, project *models.Project, tos *models.TermsOfService) error
}
