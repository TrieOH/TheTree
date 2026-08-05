package email_templates

import (
	"IdentityX/internal/emails"
	"IdentityX/models"
	"context"

	"lib/telemetry"

	"github.com/google/uuid"
)

// List returns the effective templates for both kinds — the project's
// override when present, otherwise the baked-in default — each tagged with
// its source.
func (o *Operations) List(ctx context.Context, projectID uuid.UUID) ([]models.EffectiveEmailTemplate, error) {
	ctx, span := telemetry.StartSpan(ctx, "ListEmailTemplates")
	defer span.End()

	err := o.authorizeAdmin(ctx, projectID)
	if err != nil {
		return nil, err
	}

	overrides, err := o.templates.ListByProject(ctx, projectID)
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
