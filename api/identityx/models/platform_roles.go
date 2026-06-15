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

type PlatformRoleRelation struct {
	ActorID   uuid.UUID        `json:"actor_id"`
	Role      PlatformRole     `json:"role"`
	Metadata  *json.RawMessage `json:"metadata"`
	CreatedAt time.Time        `json:"created_at"`
}
