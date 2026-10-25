// Package email_templates implements the project-scoped email template
// override endpoints (list/get/put/delete).
package handlers

import (
	"IdentityX/internal/features/email_templates/service"
	"IdentityX/internal/openapi"
	"IdentityX/models"
)

type Handlers struct {
	ops *service.Operations
}

func New(ops *service.Operations) *Handlers { return &Handlers{ops: ops} }

const module = "IdentityX"

func toEffective(t models.EffectiveEmailTemplate) openapi.EffectiveEmailTemplate {
	return openapi.EffectiveEmailTemplate{
		Kind:    openapi.EmailTemplateKind(t.Kind),
		Subject: t.Subject,
		Body:    t.Body,
		Source:  openapi.EffectiveEmailTemplateSource(t.Source),
	}
}
