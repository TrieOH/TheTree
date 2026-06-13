package globals

import "sync/atomic"

var setupComplete atomic.Bool

func SetupComplete() bool { return setupComplete.Load() }
func MarkSetupComplete()  { setupComplete.Store(true) }
