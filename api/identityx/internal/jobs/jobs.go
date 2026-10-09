// Package jobs holds IdentityX's background work: the handlers, and the one
// declaration of which jobs exist and how often the periodic ones run. The
// carrier (in-process in hosted mode, SQS plus a scheduler in managed mode)
// is chosen by the app wiring; nothing here knows it.
package jobs

import (
	"time"

	"IdentityX/internal/keys"
	"IdentityX/internal/repos"
	"IdentityX/internal/sqlc"
	"IdentityX/internal/tokens"
	"lib/email"
	libjobs "lib/jobs"
)

// Deps are what the handlers need.
type Deps struct {
	Queries      *sqlc.Queries
	ActionTokens *tokens.ActionTokenManager
	Keys         *keys.Manager
	Email        email.Sender
}

// Registry registers every IdentityX job kind.
func Registry(d Deps) *libjobs.Registry {
	reg := libjobs.NewRegistry()
	libjobs.Register(reg, NewCleanupBlacklistWorker(d.Queries).Work)
	libjobs.Register(reg, NewCleanupActionTokensWorker(d.ActionTokens).Work)
	libjobs.Register(reg, NewRotateKeysWorker(d.Keys).Work)
	libjobs.Register(reg, NewSendAuthEmailWorker(d.Email, repos.NewEmailTemplates(d.Queries)).Work)
	libjobs.Register(reg, NewSendTosUpdateWorker(d.Email, repos.NewActors(d.Queries), repos.NewEmailTemplates(d.Queries)).Work)
	return reg
}

// CleanupInterval is the blacklist and action-token sweep cadence.
const CleanupInterval = 5 * time.Minute

// Periodic declares the recurring jobs. Hosted mode runs them on tickers;
// managed mode gets the same list from `identityx schedules` and turns it
// into scheduler entries, so the cadence is declared once.
//
// Boot provisions every scope's keys inline, so the rotation job does not
// run on start: it is the ongoing heartbeat that rotates and sweeps keys
// while the service stays up.
func Periodic(rotateKeysEvery time.Duration) []libjobs.Periodic {
	return []libjobs.Periodic{
		{Job: CleanupBlacklistArgs{}, Every: CleanupInterval, RunOnStart: true},
		{Job: CleanupActionTokensArgs{}, Every: CleanupInterval, RunOnStart: true},
		{Job: RotateKeysArgs{}, Every: rotateKeysEvery, RunOnStart: false},
	}
}
