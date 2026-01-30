package scopes

import (
	"time"

	"github.com/google/uuid"
)

type ScopeType string

const (
	ScopeTypeGlobal       ScopeType = "global"
	ScopeTypeProjectRoot  ScopeType = "project_root"
	ScopeTypeProjectScope ScopeType = "project_scope"
	ScopeTypeNone         ScopeType = "none"
)

type Scope struct {
	ID         uuid.UUID
	Type       ScopeType
	ProjectID  *uuid.UUID
	Name       *string
	ExternalID *string
	CreatedAt  time.Time
}
