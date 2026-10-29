// Package email_templates is the email templates feature: its repo and
// handlers live here, its operations in service/, and Build wires the
// feature from here.
//
// The handlers serve the project-scoped email template overrides
// (list/get/put/delete).
package email_templates

import (
	"IdentityX/internal/features/email_templates/service"
)

type (
	Deps       = service.Deps
	Operations = service.Operations
)

func New(d Deps) *Operations { return service.New(d) }
