// Package repos aggregates every feature's repository layer. Import this
// package instead of the per-feature subpackages:
//
//	r := repos.New(q)
package repos

import (
	"Informd/internal/sqlc"

	"Informd/internal/repos/answers"
	"Informd/internal/repos/fields"
	"Informd/internal/repos/forms"
	"Informd/internal/repos/namespaces"
	"Informd/internal/repos/responders"
	"Informd/internal/repos/responses"
	"Informd/internal/repos/steps"
)

// Type and constructor aliases for each feature's repo package.
type (
	Namespaces = namespaces.Repo
	Forms      = forms.Repo
	Steps      = steps.Repo
	Fields     = fields.Repo
	Answers    = answers.Repo
	Responders = responders.Repo
	Responses  = responses.Repo
)

var (
	NewNamespaces = namespaces.NewRepo
	NewForms      = forms.NewRepo
	NewSteps      = steps.NewRepo
	NewFields     = fields.NewRepo
	NewAnswers    = answers.NewRepo
	NewResponders = responders.NewRepo
	NewResponses  = responses.NewRepo
)

// Repos is the aggregate of every feature repo, constructed once at startup.
type Repos struct {
	Namespaces *Namespaces
	Forms      *Forms
	Steps      *Steps
	Fields     *Fields
	Answers    *Answers
	Responders *Responders
	Responses  *Responses
}

// New constructs every feature repo from the shared query handle.
func New(q *sqlc.Queries) *Repos {
	return &Repos{
		Namespaces: NewNamespaces(q),
		Forms:      NewForms(q),
		Steps:      NewSteps(q),
		Fields:     NewFields(q),
		Answers:    NewAnswers(q),
		Responders: NewResponders(q),
		Responses:  NewResponses(q),
	}
}
