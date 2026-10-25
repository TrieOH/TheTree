package app

import (
	"net/http"
	"strings"
	"testing"

	"lib/database"
)

// TestConstraintErrorsAreMapped — every constraint's named error answers
// with its own status: a unique violation (uniq_* or a *_key index) is a
// 409, a check a 400.
func TestConstraintErrorsAreMapped(t *testing.T) {
	m := Errors()
	for name, named := range Constraints() {
		class, want := database.ErrCheckViolation, http.StatusBadRequest
		if strings.HasPrefix(name, "uniq_") || strings.HasSuffix(name, "_key") {
			class, want = database.ErrUniqueViolation, http.StatusConflict
		}
		got := m.Status(&database.ViolationError{Named: named, Class: class})
		if got != want {
			t.Errorf("%s (%v): got %d, want %d", name, named, got, want)
		}
	}
}
