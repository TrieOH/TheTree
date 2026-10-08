// Package setup holds IdentityX's setup state: whether the initial platform
// setup (/auth/setup) has completed. The database is the source of truth —
// setup is complete once any actor exists — because several processes
// (function instances, replicas) serve the same database and each must see
// a setup another one completed. Completion never reverts, so once observed
// it is cached and later checks cost nothing.
package setup

import (
	"context"
	"sync/atomic"
)

var (
	complete atomic.Bool
	probe    atomic.Pointer[func(context.Context) (bool, error)]
)

// UseProbe sets how Check asks the database whether setup has completed.
func UseProbe(fn func(context.Context) (bool, error)) { probe.Store(&fn) }

// Check reports whether setup has completed, asking the database until it
// has. With no probe configured only MarkComplete can complete it.
func Check(ctx context.Context) (bool, error) {
	if complete.Load() {
		return true, nil
	}
	fn := probe.Load()
	if fn == nil {
		return false, nil
	}
	done, err := (*fn)(ctx)
	if err != nil {
		return false, err
	}
	if done {
		complete.Store(true)
	}
	return done, nil
}

// MarkComplete records completion observed by this process.
func MarkComplete() { complete.Store(true) }
