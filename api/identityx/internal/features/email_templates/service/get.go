package service

import (
	"IdentityX/internal/emails"
	"IdentityX/models"
	"context"
	"errors"
	"lib/database"

	"github.com/google/uuid"
)

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
