package authz

import (
	"testing"

	"github.com/MintzyG/fun"
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

func isForbidden(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	var appErr *fun.AppError
	require.ErrorAs(t, err, &appErr, "expected *fun.AppError, got %T", err)
	require.Equal(t, fun.CodeForbidden, appErr.Code)
}

func TestMin(t *testing.T) {
	t.Run("below threshold is forbidden", func(t *testing.T) {
		isForbidden(t, Min(testRoleViewer, testRoleAdmin))
	})

	t.Run("equal to threshold passes", func(t *testing.T) {
		require.NoError(t, Min(testRoleAdmin, testRoleAdmin))
	})

	t.Run("above threshold passes", func(t *testing.T) {
		require.NoError(t, Min(testRoleOwner, testRoleAdmin))
	})

	t.Run("unknown role fails closed", func(t *testing.T) {
		isForbidden(t, Min(testRoleUnknown, testRoleAdmin))
	})
}

func TestForbiddenIfNotFound(t *testing.T) {
	t.Run("nil passes through", func(t *testing.T) {
		require.NoError(t, ForbiddenIfNotFound(nil))
	})

	t.Run("not found becomes forbidden", func(t *testing.T) {
		isForbidden(t, ForbiddenIfNotFound(fun.ErrNotFound("event not found")))
	})

	t.Run("other errors pass through", func(t *testing.T) {
		internal := fun.ErrInternal("boom")
		require.Equal(t, internal, ForbiddenIfNotFound(internal))
	})
}
