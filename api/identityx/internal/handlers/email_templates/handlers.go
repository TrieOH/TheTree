// Package email_templates implements the project-scoped email template
// override endpoints (list/get/put/delete).
package email_templates

import (
	"IdentityX/internal/openapi"
	"IdentityX/internal/services"
	"IdentityX/models"
)

type Handlers struct {
	ops *services.EmailTemplates
}

func New(ops *services.EmailTemplates) *Handlers { return &Handlers{ops: ops} }

const module = "IdentityX"

func toEffective(t models.EffectiveEmailTemplate) openapi.EffectiveEmailTemplate {
	return openapi.EffectiveEmailTemplate{
		Kind:    openapi.EmailTemplateKind(t.Kind),
		Subject: t.Subject,
		Body:    t.Body,
		Source:  openapi.EffectiveEmailTemplateSource(t.Source),
	}
}
