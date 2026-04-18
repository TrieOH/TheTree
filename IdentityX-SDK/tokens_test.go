package idx

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenService_ValidateToken(t *testing.T) {
	// 1. Generate Ed25519 key pair
	pub, priv, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	kid := "test-kid"
	projectID := uuid.New()

	// 2. Create a mock server for JWKS
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.URL.Path != fmt.Sprintf("/projects/%s/.well-known/jwks.json", projectID) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		jwks := JWKS{
			Keys: []JWK{
				{
					Kty: "OKP",
					Crv: "Ed25519",
					X:   base64.RawURLEncoding.EncodeToString(pub),
					Alg: "EdDSA",
					Use: "sig",
					Kid: kid,
				},
			},
		}
		json.NewEncoder(w).Encode(jwks)
	}))
	defer ts.Close()

	// 3. Initialize Client
	client, err := NewClient(Config{
		BaseURL:   ts.URL,
		APIKey:    "test-api-key",
		ProjectID: projectID,
	})
	require.NoError(t, err)

	// Verify GetJWKS manual call
	jwks, err := client.Tokens.GetJWKS(context.Background(), false)
	require.NoError(t, err)
	assert.Len(t, jwks.Keys, 1)

	// 4. Create a signed token
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"sub": "test-user",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	token.Header["kid"] = kid

	tokenStr, err := token.SignedString(priv)
	require.NoError(t, err)

	// 5. Validate the token
	validatedToken, err := client.Tokens.ValidateToken(context.Background(), tokenStr)
	require.NoError(t, err)
	assert.True(t, validatedToken.Valid)
	assert.Equal(t, "test-user", validatedToken.Claims.(jwt.MapClaims)["sub"])
}

func TestUserService_List(t *testing.T) {
	projectID := uuid.New()

	// 1. Create a mock server for List Users
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.URL.Path != fmt.Sprintf("/projects/%s/users", projectID) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		users := []ProjectUser{
			{
				ID:        uuid.New(),
				Email:     "user1@example.com",
				ProjectID: projectID,
			},
			{
				ID:        uuid.New(),
				Email:     "user2@example.com",
				ProjectID: projectID,
			},
		}
		res := struct {
			Data []ProjectUser `json:"data"`
		}{Data: users}
		json.NewEncoder(w).Encode(res)
	}))
	defer ts.Close()

	// 2. Initialize Client
	client, err := NewClient(Config{
		BaseURL:   ts.URL,
		APIKey:    "test-api-key",
		ProjectID: projectID,
	})
	require.NoError(t, err)

	// 3. List users
	users, err := client.Users.List(context.Background())
	require.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "user1@example.com", users[0].Email)
	assert.Equal(t, "user2@example.com", users[1].Email)
}

func TestUserService_Get(t *testing.T) {
	projectID := uuid.New()
	userID := uuid.New()

	// 1. Create a mock server for Get User
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.URL.Path != fmt.Sprintf("/projects/%s/users/%s", projectID, userID) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		user := ProjectUser{
			ID:        userID,
			Email:     "user1@example.com",
			ProjectID: projectID,
		}
		res := struct {
			Data ProjectUser `json:"data"`
		}{Data: user}
		json.NewEncoder(w).Encode(res)
	}))
	defer ts.Close()

	// 2. Initialize Client
	client, err := NewClient(Config{
		BaseURL:   ts.URL,
		APIKey:    "test-api-key",
		ProjectID: projectID,
	})
	require.NoError(t, err)

	// 3. Get user
	user, err := client.Users.Get(context.Background(), userID)
	require.NoError(t, err)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "user1@example.com", user.Email)
}

func TestTokenService_KeyRotation(t *testing.T) {
	// 1. Setup mock server with dynamic keys
	projectID := uuid.New()
	pub1, priv1, _ := ed25519.GenerateKey(nil)
	pub2, priv2, _ := ed25519.GenerateKey(nil)

	currentKid := "kid-1"
	currentPub := pub1

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwks := JWKS{
			Keys: []JWK{
				{
					Kty: "OKP",
					Crv: "Ed25519",
					X:   base64.RawURLEncoding.EncodeToString(currentPub),
					Alg: "EdDSA",
					Use: "sig",
					Kid: currentKid,
				},
			},
		}
		json.NewEncoder(w).Encode(jwks)
	}))
	defer ts.Close()

	client, _ := NewClient(Config{
		BaseURL:   ts.URL,
		APIKey:    "test",
		ProjectID: projectID,
	})

	// 2. Validate token with kid-1 (populates cache)
	token1 := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{"sub": "u1"})
	token1.Header["kid"] = "kid-1"
	t1Str, _ := token1.SignedString(priv1)

	_, err := client.Tokens.ValidateToken(context.Background(), t1Str)
	require.NoError(t, err)

	// 3. Rotate keys on server
	currentKid = "kid-2"
	currentPub = pub2

	// Bypass cooldown for testing
	client.Tokens.mu.Lock()
	client.Tokens.lastUpdated = time.Time{}
	client.Tokens.mu.Unlock()

	// 4. Validate token with kid-2
	// This should trigger a refresh because kid-2 is not in cache
	token2 := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{"sub": "u2"})
	token2.Header["kid"] = "kid-2"
	t2Str, _ := token2.SignedString(priv2)

	validated, err := client.Tokens.ValidateToken(context.Background(), t2Str)
	require.NoError(t, err, "Should refresh JWKS and validate new kid")
	assert.Equal(t, "u2", validated.Claims.(jwt.MapClaims)["sub"])
}
