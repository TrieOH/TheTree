// Package setup holds IdentityX's setup state: whether the initial platform
// setup (/auth/setup) has completed. The database is the source of truth —
// setup is complete once any actor exists — because several processes
// (container instances, replicas) serve the same database and each must see
// a setup another one completed. Completion never reverts, so once observed
// it is cached and later checks cost nothing.
//
// A State is built once at boot and passed to what reads it (the setup
// guard and the /auth/setup handlers); there is no package-level instance.
package setup

import (
	"context"
	"sync/atomic"
)

// State tracks setup completion for one process.
type State struct {
	complete atomic.Bool
	probe    func(context.Context) (bool, error)
}

// NewState returns a State that asks probe whether setup has completed. A
// nil probe means only MarkComplete can complete it.
func NewState(probe func(context.Context) (bool, error)) *State {
	return &State{probe: probe}
}

// Check reports whether setup has completed, asking the database until it
// has.
func (s *State) Check(ctx context.Context) (bool, error) {
	if s.complete.Load() {
		return true, nil
	}
	if s.probe == nil {
		return false, nil
	}
	done, err := s.probe(ctx)
	if err != nil {
		return false, err
	}
	if done {
		s.complete.Store(true)
	}
	return done, nil
}

// MarkComplete records completion observed by this process.
func (s *State) MarkComplete() { s.complete.Store(true) }
