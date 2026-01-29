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

	t.Run("Schema Register", func(t *testing.T) {
		testSchemaRegister(t, suite)
	})

	t.Run("Verification", func(t *testing.T) {
		testVerification(t, suite)
	})

	t.Run("Key Rotation", func(t *testing.T) {
		testKeyRotation(t, suite)
	})

	t.Run("Scopes", func(t *testing.T) {
		testScopes(t, suite)
	})

	t.Run("Permissions", func(t *testing.T) {
		testPermissions(t, suite)
	})

	t.Run("Roles", func(t *testing.T) {
		testRoles(t, suite)
	})

	t.Run("EffectivePermissions", func(t *testing.T) {
		testEffectivePermissions(t, suite)
	})

	t.Run("PermissionMatching", func(t *testing.T) {
		testPermissionMatching(t, suite)
	})

	t.Run("ConditionMatching", func(t *testing.T) {
		testConditions(t)
	})
}
