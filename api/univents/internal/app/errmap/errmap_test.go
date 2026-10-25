package errmap

import (
	"net/http"
	"strings"
	"testing"

	"lib/database"
	"univents/internal/app"
)

// TestConstraintErrorsAreMapped — every constraint's named error answers
// with its own status: a unique index (uniq_*, idx_one_*) or an exclusion
// (excl_*) is a 409, a check a 400.
func TestConstraintErrorsAreMapped(t *testing.T) {
	m := Errors()
	for name, named := range app.Constraints() {
		class, want := database.ErrCheckViolation, http.StatusBadRequest
		switch {
		case strings.HasPrefix(name, "uniq_"), strings.HasPrefix(name, "idx_one_"):
			class, want = database.ErrUniqueViolation, http.StatusConflict
		case strings.HasPrefix(name, "excl_"):
			class, want = database.ErrExclusionViolation, http.StatusConflict
		}
		got := m.Status(&database.ViolationError{Named: named, Class: class})
		if got != want {
			t.Errorf("%s (%v): got %d, want %d", name, named, got, want)
		}
	}
}
