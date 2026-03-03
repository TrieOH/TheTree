package testing

import "testing"

// ============================================================================
// ACTUAL TESTS - Clean and readable
// ============================================================================

func TestGoAuth(t *testing.T) {
	suite := NewTestSuite(t)

	t.Run("Startup", func(t *testing.T) {
		testStartup(t, suite)
	})
}
