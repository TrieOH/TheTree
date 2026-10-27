package authz

import (
	"errors"
	"testing"

	"lib/database"

	"github.com/stretchr/testify/require"
)

type testRole string

const (
	testRoleViewer  testRole = "viewer"
	testRoleAdmin   testRole = "admin"
	testRoleOwner   testRole = "owner"
	testRoleUnknown testRole = "stale"
)

func (r testRole) Rank() int {
	switch r {
	case testRoleViewer:
		return 0
	case testRoleAdmin:
		return 1
	case testRoleOwner:
		return 2
	default:
		return 0
	}
}

func (r testRole) String() string { return string(r) }

func TestAtLeast(t *testing.T) {
	t.Run("below threshold is denied", func(t *testing.T) {
		require.ErrorIs(t, AtLeast(testRoleViewer, testRoleAdmin), ErrInsufficientRole)
	})

	t.Run("equal to threshold passes", func(t *testing.T) {
		require.NoError(t, AtLeast(testRoleAdmin, testRoleAdmin))
	})

	t.Run("above threshold passes", func(t *testing.T) {
		require.NoError(t, AtLeast(testRoleOwner, testRoleAdmin))
	})

	t.Run("unknown role fails closed", func(t *testing.T) {
		require.ErrorIs(t, AtLeast(testRoleUnknown, testRoleAdmin), ErrInsufficientRole)
	})
}

func TestDenyIfNotFound(t *testing.T) {
	t.Run("nil passes through", func(t *testing.T) {
		require.NoError(t, DenyIfNotFound(nil))
	})

	t.Run("not found is denied", func(t *testing.T) {
		missing := &database.NotFoundError{Named: errors.New("event not found")}
		require.ErrorIs(t, DenyIfNotFound(missing), ErrInsufficientRole)
	})

	t.Run("other errors pass through", func(t *testing.T) {
		internal := errors.New("boom")
		require.Equal(t, internal, DenyIfNotFound(internal))
	})
}
