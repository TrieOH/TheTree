package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PlatformRole string

const (
	PlatformRoleSupport    PlatformRole = "support"
	PlatformRoleAdmin      PlatformRole = "admin"
	PlatformRoleSuperAdmin PlatformRole = "super_admin"
)

// Rank orders platform roles for lib/authz.Min: support is the lowest tier,
// super_admin the highest.
func (r PlatformRole) Rank() int {
	switch r {
	case PlatformRoleSupport:
		return 0
	case PlatformRoleAdmin:
		return 1
	case PlatformRoleSuperAdmin:
		return 2
	default:
		return 0
	}
}

func (r PlatformRole) String() string { return string(r) }

type PlatformRoleRelation struct {
	ActorID   uuid.UUID        `json:"actor_id"`
	Role      PlatformRole     `json:"role"`
	Metadata  *json.RawMessage `json:"metadata"`
	CreatedAt time.Time        `json:"created_at"`
}
