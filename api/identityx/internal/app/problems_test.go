package app

import (
	"fmt"
	"net/http"
	"testing"

	"IdentityX/errs"
	"IdentityX/internal/problems"
	"lib/database"
)

// TestConstraintErrorsAreMapped — every constraint's named error answers
// with a declared problem of its own status: a unique violation is a
// 409, a check a 400. A constraint error with no case would be a generic
// 500.
func TestConstraintErrorsAreMapped(t *testing.T) {
	for name, named := range Constraints() {
		class, want := database.ErrCheckViolation, http.StatusBadRequest
		if len(name) >= 5 && name[:5] == "uniq_" {
			class, want = database.ErrUniqueViolation, http.StatusConflict
		}
		// The bare error catches a missing case; the violation checks the
		// switch still finds the named error inside it.
		for _, err := range []error{named, &database.ViolationError{Named: named, Class: class}} {
			p := Problem(err)
			if p == nil {
				t.Errorf("%s (%v): no problem", name, named)
			} else if p.Status != want {
				t.Errorf("%s (%v): got %d, want %d", name, named, p.Status, want)
			}
		}
	}
}

// TestMostSpecificErrorWins — a named error beats the missing row it wraps.
func TestMostSpecificErrorWins(t *testing.T) {
	m := problems.Mapper(Problem)
	err := fmt.Errorf("login: %w", fmt.Errorf("%w: %w", errs.ErrInvalidCredentials, database.ErrNotFound))
	if got := m(err).Status; got != http.StatusUnauthorized {
		t.Fatalf("got %d, want 401", got)
	}
	missing := fmt.Errorf("%w: %w", errs.ErrActorNotFound, database.ErrNotFound)
	if got := m(missing).Status; got != http.StatusNotFound {
		t.Fatalf("got %d, want 404", got)
	}
}
