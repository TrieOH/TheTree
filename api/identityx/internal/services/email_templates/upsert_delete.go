package email_templates

import (
	"context"

	"IdentityX/internal/emails"
	"IdentityX/models"

	"lib/telemetry"

	"github.com/MintzyG/fun"
	"github.com/google/uuid"
)

// Upsert saves a project's override for a kind, enforcing the template
// contract (renders, and {{.ActionURL}} appears in the body at least once)
// before anything is written.
func (o *Operations) Upsert(ctx context.Context, projectID uuid.UUID, kind models.EmailTemplateKind, subject, body string) (*models.EffectiveEmailTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "UpsertEmailTemplate")
	defer span.End()

	err := o.authorizeAdmin(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if !validKind(kind) {
		return nil, invalidKindErr()
	}

	candidate := emails.Template{
		Kind:    kind,
		Subject: subject,
		Body:    body,
		Source:  emails.SourceOverride,
	}
	err = emails.Validate(candidate)
	if err != nil {
		return nil, err
	}

	saved, err := o.templates.Upsert(ctx, models.EmailTemplate{
		ProjectID: projectID,
		Kind:      kind,
		Subject:   subject,
		Body:      body,
	})
	if err != nil {
		return nil, err
	}
	return &models.EffectiveEmailTemplate{
		Kind:    kind,
		Subject: saved.Subject,
		Body:    saved.Body,
		Source:  emails.SourceOverride,
	}, nil
}

// Delete removes the project's override for a kind, restoring the default.
// Deleting a kind that has no override is a no-op (it was already the
// default).
func (o *Operations) Delete(ctx context.Context, projectID uuid.UUID, kind models.EmailTemplateKind) error {
	ctx, span := telemetry.StartSpan(ctx, "DeleteEmailTemplate")
	defer span.End()

	err := o.authorizeAdmin(ctx, projectID)
	if err != nil {
		return err
	}
	if !validKind(kind) {
		return invalidKindErr()
	}

	err = o.templates.Delete(ctx, projectID, kind)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
		return err
	}
	return nil
}
