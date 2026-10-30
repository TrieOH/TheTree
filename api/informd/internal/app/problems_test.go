package app

import (
	"net/http"
	"strings"
	"testing"

	"lib/database"
)

// TestConstraintErrorsAreMapped — every constraint's named error answers
// with a declared problem of its own status: a unique violation (uniq_*
// or a *_key index) is a 409, a check a 400.
func TestConstraintErrorsAreMapped(t *testing.T) {
	for name, named := range Constraints() {
		class, want := database.ErrCheckViolation, http.StatusBadRequest
		if strings.HasPrefix(name, "uniq_") || strings.HasSuffix(name, "_key") {
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
