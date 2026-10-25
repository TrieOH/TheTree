package app

import (
	"net/http"
	"strings"
	"testing"

	"lib/database"
	"payssage/errs"

	"github.com/google/uuid"
)

// TestConstraintErrorsAreMapped — every constraint's named error answers
// with its own status: a unique violation is a 409, a check a 400.
func TestConstraintErrorsAreMapped(t *testing.T) {
	m := Errors()
	for name, named := range Constraints() {
		class, want := database.ErrCheckViolation, http.StatusBadRequest
		if strings.HasPrefix(name, "uniq_") {
			class, want = database.ErrUniqueViolation, http.StatusConflict
		}
		got := m.Status(&database.Violation{Named: named, Class: class})
		if got != want {
			t.Errorf("%s (%v): got %d, want %d", name, named, got, want)
		}
	}
}

// TestDataDependentStatuses — a provider's answer and a revoked seller map
// by what they carry.
func TestDataDependentStatuses(t *testing.T) {
	m := Errors()
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
		if got := m.Status(tc.err); got != tc.want {
			t.Errorf("%s: got %d, want %d", tc.name, got, tc.want)
		}
	}
}
