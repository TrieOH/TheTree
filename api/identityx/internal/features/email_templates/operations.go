// Package email_templates is the email templates feature: its repo, handlers
// and operations. Build wires it.
//
// The handlers serve the project-scoped email template overrides
// (list/get/put/delete).
package email_templates

import (
	"IdentityX/errs"
	"IdentityX/internal/authz"
	"IdentityX/internal/emails"
	"IdentityX/models"
	"IdentityX/ports"
	"context"
	"errors"
	"lib/database"
	"uuid"
)

type Deps struct {
	Templates ports.EmailTemplateRepo
	Authz     *authz.Service
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }

// authorizeAdmin gates a write on the authenticated actor holding an
// admin or owner role on the project; unknown projects surface as 404
// (CheckProject passes the project lookup through).
func (o *Operations) authorizeAdmin(ctx context.Context, projectID uuid.UUID) error {
	return o.Authz.CheckProject(ctx, projectID, models.ProjectRoleAdmin)
}

func validKind(kind models.EmailTemplateKind) bool {
	return kind == models.VerifyEmailTemplateKind || kind == models.ResetEmailTemplateKind || kind == models.TosEmailTemplateKind
}

func invalidKindErr() error {
	return errs.ErrInvalidTemplateKind
}

// Get returns the effective template for a kind: the project's override
// when present, otherwise the baked-in default.
func (o *Operations) Get(ctx context.Context, projectID uuid.UUID, kind models.EmailTemplateKind) (*models.EffectiveEmailTemplate, error) {
	err := o.authorizeAdmin(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if !validKind(kind) {
		return nil, invalidKindErr()
	}

	override, err := o.Templates.GetByProjectAndKind(ctx, projectID, kind)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return nil, err
	}
	if err == nil {
		return &models.EffectiveEmailTemplate{
			Kind:    kind,
			Subject: override.Subject,
			Body:    override.Body,
			Source:  emails.SourceOverride,
		}, nil
	}

	def := emails.Default(kind)
	return &models.EffectiveEmailTemplate{
		Kind:    kind,
		Subject: def.Subject,
		Body:    def.Body,
		Source:  emails.SourceDefault,
	}, nil
}

// List returns the effective templates for both kinds — the project's
// override when present, otherwise the baked-in default — each tagged with
// its source.
func (o *Operations) List(ctx context.Context, projectID uuid.UUID) ([]models.EffectiveEmailTemplate, error) {
	err := o.authorizeAdmin(ctx, projectID)
	if err != nil {
		return nil, err
	}

	overrides, err := o.Templates.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	byKind := make(map[models.EmailTemplateKind]models.EmailTemplate, len(overrides))
	for _, t := range overrides {
		byKind[t.Kind] = t
	}

	out := make([]models.EffectiveEmailTemplate, 0, len(models.AllEmailTemplateKinds))
	for _, kind := range models.AllEmailTemplateKinds {
		if override, ok := byKind[kind]; ok {
			out = append(out, models.EffectiveEmailTemplate{
				Kind:    kind,
				Subject: override.Subject,
				Body:    override.Body,
				Source:  emails.SourceOverride,
			})
			continue
		}
		def := emails.Default(kind)
		out = append(out, models.EffectiveEmailTemplate{
			Kind:    kind,
			Subject: def.Subject,
			Body:    def.Body,
			Source:  emails.SourceDefault,
		})
	}
	return out, nil
}

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
