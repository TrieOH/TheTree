// Package services aggregates every feature's operations layer. Import this
// package instead of the per-feature subpackages:
//
//	ops := services.NewOperations(r)
package services

import (
	"Informd/internal/authz"
	"Informd/internal/repos"
	"Informd/internal/services/fields"
	"Informd/internal/services/forms"
	"Informd/internal/services/namespaces"
	"Informd/internal/services/responses"
	"Informd/internal/services/steps"
)

// Type and constructor aliases for each feature's operations package.
type (
	Namespaces = namespaces.Operations
	Forms      = forms.Operations
	Steps      = steps.Operations
	Fields     = fields.Operations
	Responses  = responses.Operations
)

var (
	NewNamespaces = namespaces.NewOperations
	NewForms      = forms.NewOperations
	NewSteps      = steps.NewOperations
	NewFields     = fields.NewOperations
	NewResponses  = responses.NewOperations
)

// Operations is the aggregate of every feature's operations, constructed
// once at startup and consumed by the HTTP handlers.
type Operations struct {
	Namespaces *Namespaces
	Forms      *Forms
	Steps      *Steps
	Fields     *Fields
	Responses  *Responses
}

// NewOperations wires every feature's operations from the shared repos.
// Authorization arrives by injection through the same seam — no
// service-locator globals.
func NewOperations(r *repos.Repos, authzSvc *authz.Service) *Operations {
	return &Operations{
		Namespaces: NewNamespaces(r.Namespaces, r.Forms, r.Steps, r.Fields, r.Answers, r.Responses, r.Responders, authzSvc),
		Forms:      NewForms(r.Forms, r.Steps, r.Namespaces, r.Fields, r.Answers, r.Responses, r.Responders, authzSvc),
		Steps:      NewSteps(r.Forms, r.Steps, r.Namespaces, authzSvc),
		Fields:     NewFields(r.Forms, r.Steps, r.Fields, r.Namespaces, authzSvc),
		Responses:  NewResponses(r.Responders, r.Responses, r.Answers, r.Forms),
	}
}
