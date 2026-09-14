// Package tos owns the project terms-of-service concern: versioned
// documents, the consent ledger, registration gating, and the
// change-notification flow.
package tos

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"IdentityX/models"

	"lib/telemetry"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// Get returns the project's current terms. Public: frontends render the
// document on signup and notice pages without authenticating. A project
// without terms (v0) surfaces as not-found.
func (o *Operations) Get(ctx context.Context, projectID uuid.UUID) (*models.TermsOfService, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetProjectTos")
	defer span.End()

	return o.tos.GetCurrent(ctx, projectID)
}

// Create introduces a project's first terms (v1). This is the v0 → v1
// migration the service performs on a live product: every existing user is
// notified exactly like on any later update, they simply migrate from "no
// terms" to v1, and their silence until the effective date is consent by
// continued use. Creating when terms already exist is a conflict — use
// Update.
func (o *Operations) Create(ctx context.Context, projectID uuid.UUID, content string) (*models.TermsOfService, error) {
	ctx, span := telemetry.StartSpan(ctx, "CreateProjectTos")
	defer span.End()

	err := o.authorizeAdmin(ctx, projectID)
	if err != nil {
		return nil, err
	}
	err = validContent(content)
	if err != nil {
		return nil, err
	}

	_, err = o.tos.GetCurrent(ctx, projectID)
	switch {
	case err == nil:
		return nil, fun.ErrConflict("project already has terms of service; update them instead")
	case !fun.Is(err, fun.CodeNotFound):
		return nil, err
	}

	tos, err := o.tos.Create(ctx, models.TermsOfService{
		ProjectID:   projectID,
		Version:     1,
		Content:     content,
		EffectiveAt: time.Now().Add(noticePeriod),
	})
	if err != nil {
		return nil, err
	}

	err = o.notifyUsers(ctx, projectID, tos)
	if err != nil {
		return nil, err
	}
	return tos, nil
}

// Update replaces the content and bumps the version. Every human user of
// the project is notified of the specific change (LGPD art. 9 §6): the
// email carries the new content, the effective date, the continued-use
// clause, and the account-termination path for those who disagree.
func (o *Operations) Update(ctx context.Context, projectID uuid.UUID, content string) (*models.TermsOfService, error) {
	ctx, span := telemetry.StartSpan(ctx, "UpdateProjectTos")
	defer span.End()

	err := o.authorizeAdmin(ctx, projectID)
	if err != nil {
		return nil, err
	}
	err = validContent(content)
	if err != nil {
		return nil, err
	}

	tos, err := o.tos.Update(ctx, projectID, content, time.Now().Add(noticePeriod))
	if err != nil {
		return nil, err
	}

	err = o.notifyUsers(ctx, projectID, tos)
	if err != nil {
		return nil, err
	}
	return tos, nil
}

// Accept records the caller's explicit agreement to the current version —
// the clickwrap path for existing users re-affirming after an update.
// Courts enforce clickwrap consent, not silence; the ledger row is the
// LGPD art. 8 §2 proof artifact. Accepting a version already accepted is
// idempotent.
func (o *Operations) Accept(ctx context.Context, projectID uuid.UUID) (*models.TosAcceptance, error) {
	ctx, span := telemetry.StartSpan(ctx, "AcceptProjectTos")
	defer span.End()

	ident, err := models.RequireIdentity(ctx)
	if err != nil {
		return nil, err
	}
	err = o.requireActorInProject(ctx, ident.Sub.ID, projectID)
	if err != nil {
		return nil, err
	}

	tos, err := o.tos.GetCurrent(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return o.tos.RecordAcceptance(ctx, models.TosAcceptance{
		ActorID:     ident.Sub.ID,
		ProjectID:   projectID,
		TosVersion:  tos.Version,
		ContentHash: contentHash(tos.Content),
		Source:      models.TosAcceptanceClickwrap,
	})
}

// ListAcceptances serves the consent ledger to project admins: who agreed
// to which version, when. This is the audit surface the burden of proof
// in LGPD art. 8 §2 demands.
func (o *Operations) ListAcceptances(ctx context.Context, projectID uuid.UUID) ([]models.TosAcceptanceWithActor, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListTosAcceptances")
	defer span.End()

	err := o.authorizeAdmin(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return o.tos.ListAcceptances(ctx, projectID)
}

// CurrentTos resolves a project's current terms; hasTos is false when the
// project has none (v0). The registration flows gate on it.
func (o *Operations) CurrentTos(ctx context.Context, projectID uuid.UUID) (tos *models.TermsOfService, hasTos bool, err error) {
	tos, err = o.tos.GetCurrent(ctx, projectID)
	if err != nil {
		if fun.Is(err, fun.CodeNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return tos, tos != nil, nil
}

// BindRegistration records a freshly registered actor's acceptance of the
// current terms. The caller owns gating (rejecting the registration when
// the acceptance flag is missing); this only pins the ledger row. It is a
// no-op for projects without terms.
func (o *Operations) BindRegistration(ctx context.Context, actorID, projectID uuid.UUID, tos *models.TermsOfService) error {
	if tos == nil {
		return nil
	}
	_, err := o.tos.RecordAcceptance(ctx, models.TosAcceptance{
		ActorID:     actorID,
		ProjectID:   projectID,
		TosVersion:  tos.Version,
		ContentHash: contentHash(tos.Content),
		Source:      models.TosAcceptanceClickwrap,
	})
	return err
}

// StampUse records a continued-use acceptance: called on a successful
// login, it pins the current version for the actor when the version has
// taken effect (past effective_at) and the actor has not accepted it yet.
// This turns the "continued use constitutes acceptance" clause from a
// legal fiction into a ledger row — weaker evidence than clickwrap, but
// documented, which is what LGPD art. 8 §2 asks for. No-ops when the
// project has no terms, before the effective date, or when the actor
// already accepted the current version.
func (o *Operations) StampUse(ctx context.Context, actorID, projectID uuid.UUID) error {
	ctx, span := telemetry.StartSpan(ctx, "StampTosUse")
	defer span.End()

	tos, hasTos, err := o.CurrentTos(ctx, projectID)
	if err != nil || !hasTos {
		return err
	}
	// The announced version is not binding until it takes effect; logins
	// inside the notice window stay in the old regime.
	if time.Now().Before(tos.EffectiveAt) {
		return nil
	}
	latest, err := o.tos.LatestAcceptedVersion(ctx, actorID, projectID)
	if err != nil {
		return err
	}
	if latest >= tos.Version {
		return nil
	}
	_, err = o.tos.RecordAcceptance(ctx, models.TosAcceptance{
		ActorID:     actorID,
		ProjectID:   projectID,
		TosVersion:  tos.Version,
		ContentHash: contentHash(tos.Content),
		Source:      models.TosAcceptanceContinuedUse,
	})
	return err
}

// notifyUsers enqueues the change-notification fan-out. A failure to
// enqueue fails the request loudly: an unannounced version bump would
// silently bind users who never got their LGPD notice.
func (o *Operations) notifyUsers(ctx context.Context, projectID uuid.UUID, tos *models.TermsOfService) error {
	project, err := o.projects.GetByID(ctx, projectID)
	if err != nil {
		return err
	}
	return o.notifier.NotifyUpdate(ctx, project, tos)
}

func validContent(content string) error {
	if len(content) == 0 {
		return fun.ErrValidation("terms of service content must not be empty")
	}
	return nil
}

// contentHash pins the SHA-256 of the accepted content, so the consent
// row proves which document the actor agreed to even after the live
// document moves on to a newer version.
func contentHash(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}
