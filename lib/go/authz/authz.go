// Package authz holds the role checks every Backend shares. The auth chains
// the spec derives (security, x-scope, x-cache) live in lib/httpserver.
package authz

import (
	"errors"

	"lib/database"
)

// Role is implemented by every member-role enum in the services.
type Role interface {
	Rank() int
	String() string
}

// ErrInsufficientRole is what the named-error role checks return (ADR-0016);
// lib/httpserver maps it to 403 for every Backend.
var ErrInsufficientRole = errors.New("your role here doesn't allow this operation; ask an admin or the owner for a higher role")

// AtLeast returns ErrInsufficientRole when role ranks below minRole.
func AtLeast(role, minRole Role) error {
	if role.Rank() < minRole.Rank() {
		return ErrInsufficientRole
	}
	return nil
}

// DenyIfNotFound turns a missing membership (database.ErrNotFound) into
// ErrInsufficientRole and passes everything else through: missing entities
// and non-members are both denied.
func DenyIfNotFound(err error) error {
	if errors.Is(err, database.ErrNotFound) {
		return ErrInsufficientRole
	}
	return err
}
