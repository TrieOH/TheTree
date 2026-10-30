package app

import (
	"net/http"
	"strings"
	"testing"

	"lib/database"
	"payssage/errs"
	"payssage/internal/problems"

	"github.com/google/uuid"
)

// TestConstraintErrorsAreMapped — every constraint's named error answers
// with a declared problem of its own status: a unique violation is a
// 409, a check a 400.
func TestConstraintErrorsAreMapped(t *testing.T) {
	for name, named := range Constraints() {
		class, want := database.ErrCheckViolation, http.StatusBadRequest
		if strings.HasPrefix(name, "uniq_") {
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

// TestDataDependentStatuses — a provider's answer and a revoked seller map
// by what they carry.
func TestDataDependentStatuses(t *testing.T) {
	m := problems.Mapper(Problem)
	cases := []struct {
		name string
		err  error
		want int
	}{
		{"provider outage", &errs.ProviderError{Provider: "mercadopago", Status: 503}, http.StatusBadGateway},
		{"provider rejects", &errs.ProviderError{Provider: "mercadopago", Status: 400}, http.StatusBadRequest},
		{"provider auth", &errs.ProviderError{Provider: "mercadopago", Status: 403}, http.StatusUnauthorized},
		{"revoked at checkout", &errs.SellerRevokedError{SellerID: uuid.New(), Op: "checkout"}, http.StatusForbidden},
		{"revoked at refund", &errs.SellerRevokedError{SellerID: uuid.New(), Op: "refund"}, http.StatusConflict},
	}
	for _, tc := range cases {
		if got := m(tc.err).Status; got != tc.want {
			t.Errorf("%s: got %d, want %d", tc.name, got, tc.want)
		}
	}
}
