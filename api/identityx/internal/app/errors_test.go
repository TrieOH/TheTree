package app

import (
	"fmt"
	"net/http"
	"testing"

	"IdentityX/errs"
	"lib/database"
)

// TestConstraintErrorsAreMapped — every constraint's named error answers
// with its own status: a unique violation is a 409, a check a 400. A
// constraint error missing from the map would fall back to its class and
// lose its message.
func TestConstraintErrorsAreMapped(t *testing.T) {
	m := Errors()
	for name, named := range Constraints() {
		var class error = database.ErrCheckViolation
		want := http.StatusBadRequest
		if len(name) >= 5 && name[:5] == "uniq_" {
			class, want = database.ErrUniqueViolation, http.StatusConflict
		}
		got := m.Status(fmt.Errorf("%w: %w", named, class))
		if got != want {
			t.Errorf("%s (%v): got %d, want %d", name, named, got, want)
		}
	}
}

// TestMostSpecificErrorWins — a named error beats the missing row it wraps.
func TestMostSpecificErrorWins(t *testing.T) {
	m := Errors()
	err := fmt.Errorf("login: %w", fmt.Errorf("%w: %w", errs.ErrInvalidCredentials, database.ErrNotFound))
	if got := m.Status(err); got != http.StatusUnauthorized {
		t.Fatalf("got %d, want 401", got)
	}
	missing := fmt.Errorf("%w: %w", errs.ErrActorNotFound, database.ErrNotFound)
	if got := m.Status(missing); got != http.StatusNotFound {
		t.Fatalf("got %d, want 404", got)
	}
}
