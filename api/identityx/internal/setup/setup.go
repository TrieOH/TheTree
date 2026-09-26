// Package setup holds IdentityX's setup-state flag: whether the initial
// platform setup (/auth/setup) has completed. The setup guard middleware
// gates every operation until it has; the setup handler marks it.
package setup

import "sync/atomic"

var complete atomic.Bool

func Complete() bool { return complete.Load() }
func MarkComplete()  { complete.Store(true) }
