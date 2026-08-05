package email_templates

import (
	"IdentityX/internal/emails"
	"IdentityX/models"
	"context"

	"lib/telemetry"

	"github.com/google/uuid"

	"github.com/MintzyG/fun"
)

// Get returns the effective template for a kind: the project's override
// when present, otherwise the baked-in default.
func (o *Operations) Get(ctx context.Context, projectID uuid.UUID, kind models.EmailTemplateKind) (*models.EffectiveEmailTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "GetEmailTemplate")
	defer span.End()

	err := o.authorizeAdmin(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if !validKind(kind) {
		return nil, invalidKindErr()
	}

	override, err := o.templates.GetByProjectAndKind(ctx, projectID, kind)
	if err != nil && !fun.Is(err, fun.CodeNotFound) {
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
