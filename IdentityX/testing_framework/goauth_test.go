package testing

import (
	"testing"
)

// ============================================================================
// ACTUAL TESTS - Clean and readable
// ============================================================================

func TestGoAuth(t *testing.T) {
	suite := NewTestSuite(t)

	t.Run("Create", func(t *testing.T) {
		testRegister(t, suite)
	})

	t.Run("Login", func(t *testing.T) {
		testLogin(t, suite)
	})

	t.Run("Sessions", func(t *testing.T) {
		testSessions(t, suite)
	})

	t.Run("Refresh", func(t *testing.T) {
		testRefresh(t, suite)
	})

	t.Run("Projects", func(t *testing.T) {
		testProjects(t, suite)
	})

	t.Run("Project Users", func(t *testing.T) {
		testProjectUsers(t, suite)
	})

	t.Run("Schemas", func(t *testing.T) {
		testSchemas(t, suite)
	})
}
