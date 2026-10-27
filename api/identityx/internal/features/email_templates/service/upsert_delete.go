package service

import (
	"context"
	"errors"

	"IdentityX/internal/emails"
	"IdentityX/models"
	"lib/database"

	"github.com/google/uuid"
)

// Upsert saves a project's override for a kind, enforcing the template
// contract (renders, and {{.ActionURL}} appears in the body at least once)
// before anything is written.
func (o *Operations) Upsert(ctx context.Context, projectID uuid.UUID, kind models.EmailTemplateKind, subject, body string) (*models.EffectiveEmailTemplate, error) {
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

	saved, err := o.Templates.Upsert(ctx, models.EmailTemplate{
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
	err := o.authorizeAdmin(ctx, projectID)
	if err != nil {
		return err
	}
	if !validKind(kind) {
		return invalidKindErr()
	}

	err = o.Templates.Delete(ctx, projectID, kind)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}
	return nil
}
