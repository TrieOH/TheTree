// Package cachetags names what the edge-cached profile reads hold
// (ADR-0020). Profiles migrate to their scope's active schema when read,
// so a schema change purges every profile in its scope as well as the
// actor's own writes purging that actor's profile.
package cachetags

import "uuid"

// ProfilesByHandle is every profile read by handle. A by-handle read does
// not know the actor's scope up front, so any schema change purges them
// all; schema changes are rare.
const ProfilesByHandle = "profiles:by-handle"

// Profile is one actor's profile, however it is addressed (by project and
// actor, or by handle, whose old URL a handle change must also clear).
func Profile(actorID uuid.UUID) string { return "profile:" + actorID.String() }

// Profiles is every profile in a scope: a project's, or the platform's
// when projectID is nil.
func Profiles(projectID *uuid.UUID) string {
	if projectID == nil {
		return "profiles:platform"
	}
	return "profiles:" + projectID.String()
}
