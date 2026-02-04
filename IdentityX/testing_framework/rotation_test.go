package testing

import (
	"GoAuth/internal/adapters/persistence/sqlc"
	"GoAuth/internal/apierr"
	"GoAuth/internal/crypto"
	"context"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/require"
)

func testKeyRotation(t *testing.T, suite *TestSuite) {
	t.Run("GlobalKeyRotation", func(t *testing.T) {
		client := suite.NewClient(t)
		user := client.WithCredentials("rotate-client@mail.com", ValidPassword).Register()
		userIsolated := client.WithCredentials("rotate-client_isolated@mail.com", ValidPassword).Register()

		// 1. Login before rotation, get KID, and store auth context
		loggedInUser := user.Login()
		preRotationKID := getKIDFromToken(t, loggedInUser.auth.AccessToken)
		preRotationAuth := loggedInUser.auth

		loggedIsolatedUser := userIsolated.Login()

		// 1.1 Verify pre-rotation token works
		client.WithAuth(preRotationAuth).GET("/sessions/me").Expect(http.StatusOK)

		// 2. Trigger global key rotation directly via sqlc
		queries := sqlc.New(suite.DB)
		ctx := context.Background()
		err := queries.RotateSigningKeysForGoAuth(ctx)
		require.NoError(t, err, "Failed to rotate GoAuth keys")

		// 2.1 Create a new active key to complete the rotation
		pub, priv, err := crypto.GenerateEd25519()
		if err != nil {
			log.Fatalf("failed to generate GoAuth key: %v", err)
		}
		defer zero(priv)

		encryptedPriv, err := crypto.Encrypt(priv)
		if err != nil {
			log.Fatalf("failed to encrypt GoAuth key: %v", err)
		}

		kid := "goauth:" + ulid.Make().String()
		expiresAt := time.Now().Add(90 * 24 * time.Hour)

		_, err = queries.CreateKeyPair(ctx, sqlc.CreateKeyPairParams{
			Kid:        kid,
			ProjectID:  nil,
			KeyType:    "goauth",
			Algorithm:  "EdDSA",
			PublicKey:  pub,
			PrivateKey: encryptedPriv,
			Usage:      "sign",
			Status:     "active",
			ExpiresAt:  expiresAt,
		})
		require.NoError(t, err, "Failed to create new active GoAuth key")

		// 3. Login after rotation and get new KID
		rotatedUser := user.Login()
		postRotationKID := getKIDFromToken(t, rotatedUser.auth.AccessToken)

		// 4. Assert KIDs are different
		require.NotEqual(t, preRotationKID, postRotationKID, "Global KID should change after rotation")

		// 5. Verify old token still works after rotation
		client.WithAuth(preRotationAuth).GET("/sessions/me").Expect(http.StatusOK)

		// 6. Refresh old token and check if it gets the new KID
		refreshedUser := client.WithAuth(preRotationAuth).Refresh()
		refreshedPreRotationKID := getKIDFromToken(t, refreshedUser.auth.AccessToken)
		require.Equal(t, postRotationKID, refreshedPreRotationKID, "Refreshed token should have the new KID")
		require.NotEqual(t, preRotationKID, refreshedPreRotationKID, "Refreshed token should have a different KID than the original")

		// --- Start Key Lifecycle Tests ---

		// 7. Manually expire the old key
		_, err = suite.DB.Exec(ctx, "UPDATE key_pair SET expires_at = $1 WHERE kid = $2", time.Now().Add(-1*time.Hour), preRotationKID)
		require.NoError(t, err, "Failed to manually expire old key")

		// 8. Verify old token still works (token's own expiry is what matters)
		client.WithAuth(loggedIsolatedUser.auth).GET("/sessions/me").Expect(http.StatusOK)

		// 9. Manually revoke the old key via the Service (Delete-Through)
		err = suite.App.Keys.RevokeKey(ctx, preRotationKID)
		require.NoError(t, err, "Failed to manually revoke old key")

		// 10. Verify old token now fails (Cache should be cleared)
		client.WithAuth(loggedIsolatedUser.auth).GET("/sessions/me").
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.TokenUntrusted).
			HasMessage("untrusted access token")

		// 11. Delete the expired and revoked key
		err = queries.DeleteExpiredRevokedKeys(ctx)
		require.NoError(t, err)

		// 12. Verify old token still fails as key is gone
		client.WithAuth(loggedIsolatedUser.auth).GET("/sessions/me").
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.TokenUntrusted).
			HasMessage("untrusted access token")
	})

	t.Run("ProjectKeyRotation", func(t *testing.T) {
		// 1. Setup project and user
		client := suite.NewClient(t)
		projectOwner := client.WithCredentials("rotate-project-owner@mail.com", ValidPassword).
			Register().
			Login().
			CreateProject("Rotation Test Project")

		projectUserClient := suite.NewClient(t)
		projectUser := projectUserClient.WithCredentials("rotate-project-user@mail.com", ValidPassword).
			ProjectRegister(projectOwner.projectID)
		projectIsolated := projectUserClient.WithCredentials("rotate-isolated-project-user@mail.com", ValidPassword).
			ProjectRegister(projectOwner.projectID)

		// 2. Login before rotation, get KID, and store auth context
		loggedInUser := projectUser.ProjectLogin(projectOwner.projectID)
		preRotationKID := getKIDFromToken(t, loggedInUser.auth.AccessToken)
		preRotationAuth := loggedInUser.auth

		loggedIsolatedUser := projectIsolated.ProjectLogin(projectOwner.projectID)

		// 2.1. Verify pre-rotation token works
		projectUserClient.WithAuth(preRotationAuth).GET("/sessions/me").Expect(http.StatusOK)

		// 3. Trigger project key rotation directly via sqlc
		queries := sqlc.New(suite.DB)
		ctx := context.Background()
		projectUUID, err := uuid.Parse(projectOwner.projectID)
		require.NoError(t, err)
		err = queries.RotateSigningKeysForProject(ctx, &projectUUID)
		require.NoError(t, err, "Failed to rotate project keys")

		// 3.1 Create a new active project key
		pub, priv, err := crypto.GenerateEd25519()
		if err != nil {
			log.Fatalf("failed to generate project key: %v", err)
		}
		defer zero(priv)

		encryptedPriv, err := crypto.Encrypt(priv)
		if err != nil {
			log.Fatalf("failed to encrypt project key: %v", err)
		}

		kid := fmt.Sprintf("project:%s:%s", projectOwner.projectID, ulid.Make().String())
		expiresAt := time.Now().Add(90 * 24 * time.Hour)

		_, err = queries.CreateKeyPair(ctx, sqlc.CreateKeyPairParams{
			Kid:        kid,
			ProjectID:  &projectUUID,
			KeyType:    "project",
			Algorithm:  "EdDSA",
			PublicKey:  pub,
			PrivateKey: encryptedPriv,
			Usage:      "sign",
			Status:     "active",
			ExpiresAt:  expiresAt,
		})
		require.NoError(t, err, "Failed to create new active project key")

		// 4. Login after rotation and get new KID
		rotatedUser := projectUser.ProjectLogin(projectOwner.projectID)
		postRotationKID := getKIDFromToken(t, rotatedUser.auth.AccessToken)

		// 5. Assert KIDs are different
		require.NotEqual(t, preRotationKID, postRotationKID, "Project KID should change after rotation")

		// 6. Verify old token still works after rotation
		projectUserClient.WithAuth(preRotationAuth).GET("/sessions/me").Expect(http.StatusOK)

		// 7. Refresh old token and check if it gets the new KID
		refreshedUser := projectUserClient.WithAuth(preRotationAuth).Refresh()
		refreshedPreRotationKID := getKIDFromToken(t, refreshedUser.auth.AccessToken)
		require.Equal(t, postRotationKID, refreshedPreRotationKID, "Refreshed token should have the new KID")
		require.NotEqual(t, preRotationKID, refreshedPreRotationKID, "Refreshed token should have a different KID than the original")

		// --- Start Key Lifecycle Tests ---

		// 8. Manually expire the old key
		_, err = suite.DB.Exec(ctx, "UPDATE key_pair SET expires_at = $1 WHERE kid = $2", time.Now().Add(-1*time.Hour), preRotationKID)
		require.NoError(t, err, "Failed to manually expire old key")

		// 9. Verify old token still works
		projectUserClient.WithAuth(loggedIsolatedUser.auth).GET("/sessions/me").Expect(http.StatusOK)

		// 10. Manually revoke the old key via the Service (Delete-Through)
		err = suite.App.Keys.RevokeKey(ctx, preRotationKID)
		require.NoError(t, err, "Failed to manually revoke old key")

		// 11. Verify old token now fails (Cache should be cleared)
		projectUserClient.WithAuth(loggedIsolatedUser.auth).GET("/sessions/me").
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.TokenUntrusted).
			HasMessage("untrusted access token")

		// 12. Delete the expired and revoked key
		err = queries.DeleteExpiredRevokedKeys(ctx)
		require.NoError(t, err)

		// 13. Verify old token still fails as key is gone
		projectUserClient.WithAuth(loggedIsolatedUser.auth).GET("/sessions/me").
			Expect(http.StatusUnauthorized).
			HasErrID(apierr.TokenUntrusted).
			HasMessage("untrusted access token")
	})
}

func getKIDFromToken(t *testing.T, tokenStr string) string {
	t.Helper()
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	require.NoError(t, err)

	kid, ok := token.Header["kid"].(string)
	require.True(t, ok, "kid not found in token header")

	return kid
}
