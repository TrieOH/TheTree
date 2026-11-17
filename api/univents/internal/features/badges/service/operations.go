package service

import (
	"context"
	"time"

	libjobs "lib/jobs"
	"univents/internal/authz"
	"univents/ports"
)

// Jobs is the background-work surface: the badge email is enqueued inside
// the caller's tx with a delay. Satisfied by the lib/jobs carriers.
type Jobs interface {
	EnqueueAfter(ctx context.Context, job libjobs.Job, delay time.Duration) error
}

type Deps struct {
	Repo          ports.BadgeTemplateRepo
	Emissions     ports.BadgeEmissionRepo
	Registrations ports.RegistrationRepo
	Editions      ports.EditionRepo
	Events        ports.EventRepo
	Actors        ports.ActorResolver
	Jobs          Jobs
	Authz         *authz.Service
	AppURL        string // the Univents front, for the profile links on badges
}

type Operations struct{ Deps }

func New(d Deps) *Operations { return &Operations{Deps: d} }
