package app

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"Informd/internal/problems"
	kit "lib/backend/identityx"
	"lib/database"
	"lib/httpserver"

	idx "sdk/identityx"
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

// TestAuthErrorsAreMapped — every error the auth step can reject a request
// with (the harness middlewares, the IdentityX SDK's token checks and the
// identityx kit's authenticators) answers with its declared problem, never
// the generic 500: the auth middleware writes them through the mapper.
func TestAuthErrorsAreMapped(t *testing.T) {
	m := problems.Mapper(Problem)
	cases := []struct {
		err  error
		want int
	}{
		{httpserver.ErrMissingBearer, http.StatusUnauthorized},
		{httpserver.ErrMissingAPIKey, http.StatusUnauthorized},
		{httpserver.ErrNoCredentials, http.StatusUnauthorized},
		{&idx.InvalidTokenError{Cause: errors.New("token is malformed")}, http.StatusUnauthorized},
		{&idx.InvalidTokenError{Cause: &idx.KeyNotFoundError{Kid: "k1"}}, http.StatusUnauthorized},
		{&idx.TokenExpiredError{ExpiredAt: time.Now()}, http.StatusUnauthorized},
		{&idx.TokenNotYetValidError{ValidAt: time.Now()}, http.StatusUnauthorized},
		{&idx.InvalidIssuerError{Got: "x", Expected: "y"}, http.StatusUnauthorized},
		{kit.ErrInvalidAPIKey, http.StatusUnauthorized},
		{kit.ErrAPIKeysNotAccepted, http.StatusUnauthorized},
		{kit.ErrEmailNotVerified, http.StatusUnauthorized},
		{fmt.Errorf("%w: %w", kit.ErrIntrospectionFailed, errors.New("dial tcp: refused")), http.StatusServiceUnavailable},
	}
	for _, tc := range cases {
		if got := m(tc.err).Status; got != tc.want {
			t.Errorf("%v: got %d, want %d", tc.err, got, tc.want)
		}
	}
}
