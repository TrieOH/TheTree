package store

import (
	"time"
)

// Test-only accessors (export_test pattern) for the keepalive interval.

// SseKeepaliveInterval returns the current keepalive interval.
func SseKeepaliveInterval() time.Duration { return time.Duration(sseKeepaliveIntervalNS.Load()) }

// SetSseKeepaliveInterval overrides the keepalive interval. Test-only.
func SetSseKeepaliveInterval(d time.Duration) { sseKeepaliveIntervalNS.Store(int64(d)) }
