package authz

import "github.com/MintzyG/fun"

// Role is implemented by every member-role enum in the services.
type Role interface {
	Rank() int
	String() string
}

// Min returns Forbidden when role ranks below minRole.
func Min(role, minRole Role) error {
	if role.Rank() < minRole.Rank() {
		return fun.ErrForbidden("insufficient permissions")
	}
	return nil
}

// ForbiddenIfNotFound maps CodeNotFound to Forbidden and passes everything
// else through — missing entities and non-members are both Forbidden.
func ForbiddenIfNotFound(err error) error {
	if err == nil {
		return nil
	}
	if fun.Is(err, fun.CodeNotFound) {
		return fun.ErrForbidden("insufficient permissions")
	}
	return err
}
