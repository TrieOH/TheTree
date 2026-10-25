// Package email_templates is the email templates feature: the layers live in the
// subpackages, and Build wires the feature from here.
package email_templates

import (
	"IdentityX/internal/features/email_templates/handlers"
	"IdentityX/internal/features/email_templates/repo"
	"IdentityX/internal/features/email_templates/service"
	"IdentityX/internal/sqlc"
)

type (
	Repo       = repo.Repo
	Deps       = service.Deps
	Operations = service.Operations
	Handlers   = handlers.Handlers
)

func NewRepo(q *sqlc.Queries) *Repo { return repo.New(q) }

func New(d Deps) *Operations { return service.New(d) }

func NewHandlers(ops *Operations) *Handlers { return handlers.New(ops) }
